// querent.js — 命主档案（强制必填）
//
// 行为：
//   - 任何页面加载时若 localStorage 无档案 → 强制弹层，不可关闭，全部必填
//   - 顶栏小标显示当前命主，点击可重设（重设时可取消）
//   - 仅本地存储（localStorage），不上传
//
// 数据 schema:
//   {
//     name, gender,
//     birthYear, birthMonth, birthDay, birthHour,
//     birthProvince, birthCity,
//     currentProvince, currentCity
//   }
//   生肖（benMing）由出生年自动推算，不入 schema
//   省 === '海外' 时，city 为用户自由输入文本
//
// 暴露：
//   window.Querent.get() -> obj | null
//   window.Querent.set(obj)
//   window.Querent.clear()
//   window.Querent.benMing()       // 由 birthYear 推算的生肖（繁体）
//   window.Querent.requireOrPrompt({force?: bool, onSet?: fn})
//   window.Querent.openEdit()      // 打开可取消的重设层
//   window.Querent.onChange(cb)

(function () {
  'use strict';

  const STORAGE_KEY = 'zhouyi_querent';
  // 出生年 → 生肖（地支次序：子鼠丑牛寅虎卯兔辰龍巳蛇午馬未羊申猴酉雞戌狗亥豬）
  // 1900 年为庚子年（鼠）。zodiacByYear(1990) = 馬。
  const ZODIAC = ['鼠','牛','虎','兔','龍','蛇','馬','羊','猴','雞','狗','豬'];
  function zodiacByYear(year) {
    if (!year || year < 1900) return '';
    const idx = ((year - 1900) % 12 + 12) % 12;
    return ZODIAC[idx];
  }

  const SHICHEN = [
    { v: '子', label: '子時（23:00–01:00）' },
    { v: '丑', label: '丑時（01:00–03:00）' },
    { v: '寅', label: '寅時（03:00–05:00）' },
    { v: '卯', label: '卯時（05:00–07:00）' },
    { v: '辰', label: '辰時（07:00–09:00）' },
    { v: '巳', label: '巳時（09:00–11:00）' },
    { v: '午', label: '午時（11:00–13:00）' },
    { v: '未', label: '未時（13:00–15:00）' },
    { v: '申', label: '申時（15:00–17:00）' },
    { v: '酉', label: '酉時（17:00–19:00）' },
    { v: '戌', label: '戌時（19:00–21:00）' },
    { v: '亥', label: '亥時（21:00–23:00）' }
  ];

  const listeners = [];

  function read() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return null;
      const obj = JSON.parse(raw);
      if (!obj || typeof obj !== 'object') return null;
      const required = ['name','gender','birthYear','birthMonth','birthDay','birthHour',
                        'birthProvince','birthCity','currentProvince','currentCity'];
      for (const k of required) {
        if (!obj[k]) return null;
      }
      return {
        name: String(obj.name).trim(),
        gender: obj.gender,
        birthYear: parseInt(obj.birthYear, 10) || 0,
        birthMonth: parseInt(obj.birthMonth, 10) || 0,
        birthDay: parseInt(obj.birthDay, 10) || 0,
        birthHour: obj.birthHour,
        birthProvince: obj.birthProvince,
        birthCity: String(obj.birthCity).trim(),
        currentProvince: obj.currentProvince,
        currentCity: String(obj.currentCity).trim()
      };
    } catch (_) { return null; }
  }

  function write(obj) {
    try {
      if (!obj) localStorage.removeItem(STORAGE_KEY);
      else localStorage.setItem(STORAGE_KEY, JSON.stringify(obj));
    } catch (_) {}
    listeners.forEach(fn => { try { fn(read()); } catch (_) {} });
    renderTag();
  }

  function tr(key, fallback) {
    if (window.I18n && window.I18n.t) return window.I18n.t(key, fallback);
    return fallback;
  }

  function summary(q) {
    if (!q) return tr('querent.tag.unset', '命主：未設');
    const z = zodiacByYear(q.birthYear);
    const parts = [q.name, q.gender, String(q.birthYear) + (z ? '(' + z + ')' : '')];
    return tr('querent.tag.label', '命主：') + parts.join('·');
  }

  // ===== 右上角固定控件组（命主小标 + 语言切换）=====
  // 创建一个固定在页面右上角的容器，把页面里的语言切换节点搬进来，
  // 命主小标放在它左侧。移动 DOM 节点会保留其原有事件绑定，故无需改 i18n.js。
  function ensureTopRight() {
    let host = document.getElementById('topRight');
    if (!host) {
      host = document.createElement('div');
      host.id = 'topRight';
      host.className = 'top-right';
      document.body.appendChild(host);
    }
    // 搬入语言切换（页面各处的 .lang-switch，通常在 nav 内）
    const lang = document.querySelector('.lang-switch');
    if (lang && lang.parentElement !== host) {
      host.appendChild(lang);
    }
    return host;
  }

  function ensureTag() {
    const host = ensureTopRight();
    let tag = document.getElementById('querentTag');
    if (tag) {
      // 命主小标始终排在语言切换左侧
      if (tag.parentElement !== host) host.insertBefore(tag, host.firstChild);
      return tag;
    }
    tag = document.createElement('a');
    tag.id = 'querentTag';
    tag.className = 'querent-tag';
    tag.href = 'javascript:void(0)';
    tag.setAttribute('role', 'button');
    tag.addEventListener('click', () => openEdit());
    host.insertBefore(tag, host.firstChild); // 命主在左、语言在右
    return tag;
  }

  function renderTag() {
    const tag = ensureTag();
    if (!tag) return;
    const q = read();
    tag.textContent = summary(q);
    tag.classList.toggle('is-set', !!q);
  }

  // ===== 省/市级联控件构建 =====
  // returns: { wrap: HTMLElement, getProvince: ()=>str, getCity: ()=>str, setValue: (p,c)=>void, onChange: (fn)=>void }
  function buildProvinceCity(idPrefix, currentProvince, currentCity) {
    const wrap = document.createElement('div');
    wrap.className = 'qm-pc';
    const provSel = document.createElement('select');
    provSel.id = idPrefix + 'Province';
    const citySel = document.createElement('select');
    citySel.id = idPrefix + 'City';
    const cityInput = document.createElement('input');
    cityInput.type = 'text';
    cityInput.id = idPrefix + 'CityInput';
    cityInput.maxLength = 64;
    cityInput.placeholder = tr('querent.field.overseascity.ph', '請輸入國家/城市');
    cityInput.style.display = 'none';
    cityInput.autocomplete = 'off';

    // 省 options
    const provs = (window.Regions && window.Regions.provinces()) || [];
    provSel.innerHTML = '<option value="">' + tr('querent.opt.province', '省／地區') + '</option>' +
      provs.map(p => `<option value="${p}">${p}</option>`).join('');

    function refreshCities(p) {
      const cs = (window.Regions && window.Regions.cities(p)) || [];
      if (p === '海外') {
        citySel.style.display = 'none';
        cityInput.style.display = '';
        return;
      }
      citySel.style.display = '';
      cityInput.style.display = 'none';
      citySel.innerHTML = '<option value="">' + tr('querent.opt.city', '市／區') + '</option>' +
        cs.map(c => `<option value="${c}">${c}</option>`).join('');
    }
    refreshCities(currentProvince || '');
    if (currentProvince) provSel.value = currentProvince;
    if (currentCity) {
      if (currentProvince === '海外') cityInput.value = currentCity;
      else citySel.value = currentCity;
    }

    let externalChange = null;
    provSel.addEventListener('change', () => {
      refreshCities(provSel.value);
      if (externalChange) externalChange();
    });
    citySel.addEventListener('change', () => { if (externalChange) externalChange(); });
    cityInput.addEventListener('input', () => { if (externalChange) externalChange(); });

    wrap.appendChild(provSel);
    wrap.appendChild(citySel);
    wrap.appendChild(cityInput);

    return {
      wrap,
      getProvince: () => provSel.value,
      getCity: () => provSel.value === '海外' ? cityInput.value.trim() : citySel.value,
      onChange: fn => { externalChange = fn; }
    };
  }

  // ===== 弹层 =====
  function openDialog(mode, onSet) {
    closeDialog();
    const isRequire = mode === 'require';
    const cur = read() || {};
    const mask = document.createElement('div');
    mask.className = 'qm-mask' + (isRequire ? ' qm-require' : '');
    mask.id = 'querentMask';

    const yearOpts = (function(){
      const now = new Date().getFullYear();
      const arr = [];
      for (let y = now; y >= 1900; y--) arr.push(`<option value="${y}">${y}</option>`);
      return arr.join('');
    })();
    const monthOpts = Array.from({length:12}, (_,i)=>`<option value="${i+1}">${i+1}</option>`).join('');
    const dayOpts = Array.from({length:31}, (_,i)=>`<option value="${i+1}">${i+1}</option>`).join('');
    const hourOpts = SHICHEN.map(s=>`<option value="${s.v}">${s.label}</option>`).join('');

    mask.innerHTML = `
      <div class="qm-dialog qm-dialog-wide" role="dialog" aria-modal="true">
        <h2>${isRequire ? tr('querent.dialog.title.require', '命　主　設　定（必填）') : tr('querent.dialog.title.edit', '命　主　設　定')}</h2>
        <p class="qm-intro">${isRequire ? tr('querent.dialog.intro.require', '占卦之前，請先以正心填錄命主資訊，以為儀式之憑。所填資料僅存於本機瀏覽器，不會上傳。') : tr('querent.dialog.intro.edit', '所填資料僅存於本機瀏覽器，不會上傳。')}</p>

        <div class="qm-row">
          <label>${tr('querent.field.name', '姓　名')}<span class="req">*</span></label>
          <input type="text" id="qmName" maxlength="32" autocomplete="off" />
        </div>

        <div class="qm-row">
          <label>${tr('querent.field.gender', '性　別')}<span class="req">*</span></label>
          <select id="qmGender">
            <option value="">${tr('querent.opt.placeholder', '— 請選 —')}</option>
            <option value="男">${tr('liuren.gender.male', '男')}</option>
            <option value="女">${tr('liuren.gender.female', '女')}</option>
          </select>
        </div>

        <div class="qm-row qm-row-3">
          <label>${tr('querent.field.birthdate', '出生公曆')}<span class="req">*</span></label>
          <div class="qm-triple">
            <select id="qmYear"><option value="">${tr('querent.opt.year', '年')}</option>${yearOpts}</select>
            <select id="qmMonth"><option value="">${tr('querent.opt.month', '月')}</option>${monthOpts}</select>
            <select id="qmDay"><option value="">${tr('querent.opt.day', '日')}</option>${dayOpts}</select>
          </div>
        </div>

        <div class="qm-row">
          <label>${tr('querent.field.hour', '出生時辰')}<span class="req">*</span></label>
          <select id="qmHour">
            <option value="">${tr('querent.opt.placeholder', '— 請選 —')}</option>
            ${hourOpts}
          </select>
        </div>

        <div class="qm-row qm-row-3">
          <label>${tr('querent.field.birthplace', '出　生　地')}<span class="req">*</span></label>
          <div class="qm-pc-slot" id="qmBirthSlot"></div>
        </div>

        <div class="qm-row qm-row-3">
          <label>${tr('querent.field.curplace', '現　居　地')}<span class="req">*</span></label>
          <div class="qm-pc-slot" id="qmCurSlot"></div>
        </div>

        <p class="qm-hint">${tr('querent.dialog.hint',
          '※ 上述資料僅存於本機瀏覽器，不會主動上傳。其中性別、出生年用於六壬「年命」推算；姓名、出生地、現居地、月日時辰<strong>不參與卦象算法</strong>，僅為起卦時之儀式化標識。')}</p>

        <div class="qm-actions">
          ${isRequire ? '' : `<button type="button" class="qm-cancel" id="qmCancel">${tr('querent.btn.cancel', '取　消')}</button>`}
          <button type="button" class="qm-save" id="qmSave" disabled>${tr('querent.btn.save', '保　存')}</button>
        </div>
      </div>`;
    document.body.appendChild(mask);
    if (window.I18n && window.I18n.apply) window.I18n.apply(mask);
    document.body.classList.add('qm-open');

    const $ = id => mask.querySelector(id);

    // 注入两组省/市级联
    const birthPC = buildProvinceCity('qmBirth', cur.birthProvince, cur.birthCity);
    const curPC = buildProvinceCity('qmCur', cur.currentProvince, cur.currentCity);
    $('#qmBirthSlot').appendChild(birthPC.wrap);
    $('#qmCurSlot').appendChild(curPC.wrap);

    // 预填其余字段
    $('#qmName').value = cur.name || '';
    $('#qmGender').value = cur.gender || '';
    $('#qmYear').value = cur.birthYear || '';
    $('#qmMonth').value = cur.birthMonth || '';
    $('#qmDay').value = cur.birthDay || '';
    $('#qmHour').value = cur.birthHour || '';

    function valid() {
      return $('#qmName').value.trim() &&
             $('#qmGender').value &&
             $('#qmYear').value &&
             $('#qmMonth').value &&
             $('#qmDay').value &&
             $('#qmHour').value &&
             birthPC.getProvince() && birthPC.getCity() &&
             curPC.getProvince() && curPC.getCity();
    }
    function refresh() { $('#qmSave').disabled = !valid(); }
    mask.querySelectorAll('input,select').forEach(el => {
      el.addEventListener('input', refresh);
      el.addEventListener('change', refresh);
    });
    birthPC.onChange(refresh);
    curPC.onChange(refresh);
    refresh();

    if (!isRequire) {
      mask.addEventListener('click', e => { if (e.target === mask) closeDialog(); });
      const esc = e => { if (e.key === 'Escape') { closeDialog(); document.removeEventListener('keydown', esc); } };
      document.addEventListener('keydown', esc);
      const cancel = $('#qmCancel');
      if (cancel) cancel.addEventListener('click', closeDialog);
    } else {
      const block = e => { if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); } };
      mask._block = block;
      document.addEventListener('keydown', block, true);
    }

    $('#qmSave').addEventListener('click', () => {
      if (!valid()) return;
      const obj = {
        name: $('#qmName').value.trim(),
        gender: $('#qmGender').value,
        birthYear: parseInt($('#qmYear').value, 10),
        birthMonth: parseInt($('#qmMonth').value, 10),
        birthDay: parseInt($('#qmDay').value, 10),
        birthHour: $('#qmHour').value,
        birthProvince: birthPC.getProvince(),
        birthCity: birthPC.getCity(),
        currentProvince: curPC.getProvince(),
        currentCity: curPC.getCity()
      };
      write(obj);
      closeDialog();
      if (typeof onSet === 'function') onSet(obj);
    });
  }

  function closeDialog() {
    const m = document.getElementById('querentMask');
    if (!m) return;
    if (m._block) document.removeEventListener('keydown', m._block, true);
    m.remove();
    document.body.classList.remove('qm-open');
  }

  function openEdit() {
    if (read()) openDialog('edit');
    else openDialog('require');
  }

  function requireOrPrompt(opts) {
    opts = opts || {};
    const q = read();
    if (q && !opts.force) {
      if (typeof opts.onSet === 'function') opts.onSet(q);
      return;
    }
    openDialog('require', opts.onSet);
  }

  function boot() {
    renderTag();
    document.addEventListener('click', function (e) {
      const t = e.target;
      if (t && t.matches && t.matches('[data-lang-switch]')) {
        setTimeout(renderTag, 30);
      }
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', boot);
  } else {
    boot();
  }

  window.Querent = {
    get: read,
    set: write,
    clear: () => write(null),
    benMing: () => { const q = read(); return q ? zodiacByYear(q.birthYear) : ''; },
    zodiacByYear,
    openEdit,
    requireOrPrompt,
    onChange: (fn) => { if (typeof fn === 'function') listeners.push(fn); }
  };
})();
