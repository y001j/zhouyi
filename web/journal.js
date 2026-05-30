// journal.js — 「筮錄」本地笔记本 + 「印心」结果卡
//
// 全部本地存储（localStorage key = 'zhouyi_journal'），不上传。
//
// 数据 schema：见 /web/journal.html 注释。
//
// 暴露：window.Journal = {
//   list(),            // 倒序列出全部
//   get(id),
//   add(record),       // 返回 id
//   update(id, patch),
//   remove(id),
//   stats(),           // { total, fulfilled, unfulfilled, partial, pending }
//   exportJSON(),
//   clearAll(),
//   renderYinXin(opts) // 在结果区顶部注入「印心」小卡，opts = { kind, data, mountEl|mountId, question, headline, snapshot }
// }

(function () {
  'use strict';

  const STORAGE_KEY = 'zhouyi_journal';
  const MAX_RECORDS = 200;

  function tr(key, fallback) {
    if (window.I18n && window.I18n.t) return window.I18n.t(key, fallback);
    return fallback;
  }

  // ---- 存取 ----
  function readAll() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return [];
      const arr = JSON.parse(raw);
      return Array.isArray(arr) ? arr : [];
    } catch (_) { return []; }
  }
  function writeAll(arr) {
    try { localStorage.setItem(STORAGE_KEY, JSON.stringify(arr)); } catch (_) {}
  }
  function genId() {
    return 'jr_' + Date.now() + '_' + Math.random().toString(36).slice(2, 6);
  }

  function list() {
    return readAll().slice().sort((a, b) => (b.timestamp || 0) - (a.timestamp || 0));
  }
  function get(id) {
    return readAll().find(r => r.id === id) || null;
  }
  function add(rec) {
    const all = readAll();
    if (all.length >= MAX_RECORDS) {
      alert(tr('journal.full', '筮錄已達 200 卦上限，請先導出舊卦或於筮錄頁清理後再封藏。'));
      return null;
    }
    rec.id = rec.id || genId();
    rec.timestamp = rec.timestamp || Date.now();
    rec.verification = rec.verification || null;
    all.push(rec);
    writeAll(all);
    return rec.id;
  }
  function update(id, patch) {
    const all = readAll();
    const i = all.findIndex(r => r.id === id);
    if (i < 0) return false;
    all[i] = Object.assign({}, all[i], patch);
    writeAll(all);
    return true;
  }
  function remove(id) {
    writeAll(readAll().filter(r => r.id !== id));
  }
  function clearAll() { writeAll([]); }

  function stats() {
    const all = readAll();
    let f = 0, u = 0, p = 0, pend = 0;
    for (const r of all) {
      if (!r.verification) pend++;
      else if (r.verification.status === 'fulfilled') f++;
      else if (r.verification.status === 'unfulfilled') u++;
      else if (r.verification.status === 'partial') p++;
      else pend++;
    }
    return { total: all.length, fulfilled: f, unfulfilled: u, partial: p, pending: pend };
  }

  function exportJSON() {
    const data = JSON.stringify(readAll(), null, 2);
    const blob = new Blob([data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    const ts = new Date().toISOString().slice(0, 10);
    a.href = url;
    a.download = 'zhouyi_journal_' + ts + '.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  // ---- 印心小卡 ----
  // opts:
  //   kind: 'zhouyi' | 'liuren' | 'qimen' | 'huican'
  //   data: 后端返回的整个 result 对象（用作 snapshot）
  //   mountEl 或 mountId: 注入位置（卡片插入到该元素的 firstChild 之前）
  //   question: 用户所问
  //   headline: 卦象/课名/局名 一行小标题
  //   ganzhiTime: 干支时刻（可选，缺则用本地时刻字符串）
  //   questionType: 类别 key（可选）
  function renderYinXin(opts) {
    const mount = opts.mountEl || (opts.mountId && document.getElementById(opts.mountId));
    if (!mount) return;
    const old = mount.querySelector(':scope > .yinxin-card');
    if (old) old.remove();

    const ts = Date.now();
    const ganzhi = opts.ganzhiTime || new Date(ts).toLocaleString('zh-CN', { hour12: false });
    const card = document.createElement('div');
    card.className = 'yinxin-card';
    card.innerHTML = `
      <div class="yx-title">${tr('yinxin.title', '印　心')}</div>
      <div class="yx-meta">
        <div class="yx-q">${tr('yinxin.youask', '你問：')}<span class="yx-qt">${escape(opts.question || tr('yinxin.noquestion', '（未明所問）'))}</span></div>
        <div class="yx-when">${tr('yinxin.atwhen', '於')} <strong>${escape(ganzhi)}</strong>${tr('yinxin.got', '，得')} <strong class="yx-head">${escape(opts.headline || '')}</strong></div>
      </div>
      <div class="yx-thought">
        <label class="yx-thought-lbl">${tr('yinxin.firstthought', '你看到此卦的第一個念頭是——')}</label>
        <textarea class="yx-thought-input" rows="2" maxlength="500" placeholder="${tr('yinxin.thoughtph', '寫下這一念。它常比後來的解卦更近本心。')}"></textarea>
        <div class="yx-thought-hint">${tr('yinxin.thoughthint', '（可空。封藏后此處仍可繼續編輯，自動保存。）')}</div>
      </div>
      <div class="yx-actions">
        <button type="button" class="yx-seal-btn">${tr('yinxin.btn.seal', '封　藏　此　卦')}</button>
        <span class="yx-status"></span>
        <a href="/journal.html" class="yx-link">${tr('yinxin.btn.viewjournal', '往　筮　錄')}</a>
      </div>
    `;
    mount.insertBefore(card, mount.firstChild);
    if (window.I18n && window.I18n.apply) window.I18n.apply(card);

    const ta = card.querySelector('.yx-thought-input');
    const sealBtn = card.querySelector('.yx-seal-btn');
    const status = card.querySelector('.yx-status');

    let recId = null;
    let saveTimer = null;

    function setSealed(id) {
      recId = id;
      sealBtn.disabled = true;
      sealBtn.textContent = tr('yinxin.sealed', '已封藏 ✓');
      sealBtn.classList.add('is-sealed');
    }

    sealBtn.addEventListener('click', () => {
      if (recId) return;
      const rec = {
        kind: opts.kind,
        timestamp: ts,
        ganzhiTime: ganzhi,
        question: opts.question || '',
        questionType: opts.questionType || '',
        headline: opts.headline || '',
        snapshot: opts.data || null,
        firstThought: ta.value.trim()
      };
      const id = add(rec);
      if (id) {
        setSealed(id);
        status.textContent = tr('yinxin.justsealed', '已封藏入筮錄');
        setTimeout(() => { status.textContent = ''; }, 2200);
      }
    });

    ta.addEventListener('input', () => {
      if (!recId) return;
      clearTimeout(saveTimer);
      saveTimer = setTimeout(() => {
        update(recId, { firstThought: ta.value.trim() });
        status.textContent = tr('yinxin.saved', '已自動保存');
        setTimeout(() => { status.textContent = ''; }, 1500);
      }, 600);
    });
  }

  // 简易转义，防止 headline / question 含 < > 破坏布局
  function escape(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
  }

  window.Journal = {
    list, get, add, update, remove, stats, exportJSON, clearAll, renderYinXin,
    MAX_RECORDS
  };
})();
