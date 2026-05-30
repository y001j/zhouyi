// 起占仪式动画（共享 / 复用周易页 ritual-overlay 体系的样式）
// 用法：
//   const res = await Ritual.run({ kind:'liuren', question:'…', fetchPromise });
// 阶段：1) 凝神倒计 3-2-1（约 2.7s）  2) 起占视觉动画（最少 2.5s，与 fetch 并发）
// fetch 完成后即进入收尾，遮罩淡出，返回 fetch 的 Response。
// 若 fetch 抛错则同样收尾后向外抛出。
(function () {
  'use strict';
  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

  // 各 preset 的 title1/title2/hint2/quote1 文案以鍵的形式存在 i18n 字典中
  // （EN 模式從字典取，中文模式從 fallback 取，再由 i18n.tr/t2s 處理繁簡）
  const PRESETS = {
    liuren: {
      seal: '壬',
      title1Key: 'ritual.meditate', title1: '凝神靜氣',
      quote1: '月將加時　四課三傳<br/>以時起課　以心觀之',
      title2Key: 'ritual.liuren.title2', title2: '月　將　加　時',
      hint2Key: 'ritual.liuren.hint2', hint2: '十二支位　次第輪轉……',
      visual: 'wheel',
    },
    qimen: {
      seal: '奇',
      title1Key: 'ritual.meditate', title1: '凝神靜氣',
      quote1: '陰陽順逆　神藏鬼沒<br/>九宮八門　三奇六儀',
      title2Key: 'ritual.qimen.title2', title2: '布　局　起　盤',
      hint2Key: 'ritual.qimen.hint2', hint2: '三奇六儀　次第入宮……',
      visual: 'palaces',
    },
    huican: {
      seal: '參',
      title1Key: 'ritual.meditate', title1: '凝神靜氣',
      quote1: '一問三占　互為印證<br/>卦顯象義　壬主人事　奇定機方',
      title2Key: 'ritual.huican.title2', title2: '三　式　同　起',
      hint2Key: 'ritual.huican.hint2', hint2: '錢　·　將　·　宮　三盤齊發……',
      visual: 'trio',
    },
  };
  function L(key, fallback) {
    return (window.I18n && window.I18n.t) ? window.I18n.t(key, fallback) : fallback;
  }

  function ensureOverlay() {
    let o = document.getElementById('ritualOverlay');
    if (o) return o;
    o = document.createElement('div');
    o.id = 'ritualOverlay';
    o.className = 'ritual-overlay';
    o.hidden = true;
    o.innerHTML = `
      <div class="ritual-bg"></div>
      <div class="ritual-content">
        <div class="ritual-seal" id="ritualSeal">易</div>
        <div class="ritual-stage" id="rStageMeditate" hidden>
          <div class="ritual-title" id="rTitle1">凝神靜氣</div>
          <div class="ritual-quote" id="rQuote1"></div>
          <div class="ritual-question" id="rQuestion"></div>
          <div class="ritual-countdown" id="rCountdown">三</div>
          <div class="ritual-hint">默念所問之事……</div>
        </div>
        <div class="ritual-stage" id="rStageVisual" hidden>
          <div class="ritual-title" id="rTitle2"></div>
          <div class="ritual-visual" id="rVisual"></div>
          <div class="ritual-hint" id="rHint2"></div>
        </div>
        <button id="rSkip" class="ritual-skip" data-i18n="ritual.skip">跳過儀式</button>
      </div>
    `;
    document.body.appendChild(o);
    return o;
  }

  function show(o) {
    o.hidden = false;
    document.body.style.overflow = 'hidden';
  }
  function hide(o) {
    o.hidden = true;
    document.body.style.overflow = '';
  }

  // ==== 各类视觉 ====
  function renderVisual(container, kind) {
    container.innerHTML = '';
    container.className = 'ritual-visual';
    if (kind === 'wheel') {
      // 大六壬：十二地支圆盘
      container.classList.add('wheel');
      const zhi = ['子','丑','寅','卯','辰','巳','午','未','申','酉','戌','亥'];
      const inner = document.createElement('div');
      inner.className = 'wheel-disc';
      zhi.forEach((z, i) => {
        const cell = document.createElement('div');
        cell.className = 'wheel-cell';
        cell.style.transform = `rotate(${i * 30}deg) translateY(-72px) rotate(${-i * 30}deg)`;
        cell.style.animationDelay = (i * 0.06) + 's';
        cell.textContent = z;
        inner.appendChild(cell);
      });
      const center = document.createElement('div');
      center.className = 'wheel-center';
      center.textContent = '將';
      inner.appendChild(center);
      container.appendChild(inner);
    } else if (kind === 'palaces') {
      // 奇门：九宫
      container.classList.add('palaces');
      const labels = ['四','九','二','三','五','七','八','一','六']; // 后天九宫数
      const stems = ['乙','丙','丁','戊','己','庚','辛','壬','癸'];
      labels.forEach((n, i) => {
        const cell = document.createElement('div');
        cell.className = 'palace-cell';
        if (i === 4) cell.classList.add('center');
        cell.style.animationDelay = (i * 0.08) + 's';
        cell.innerHTML = `<span class="p-num">${n}</span><span class="p-stem">${stems[i]}</span>`;
        container.appendChild(cell);
      });
    } else if (kind === 'trio') {
      // 三式互参：左铜钱 / 中转盘 / 右九宫
      container.classList.add('trio');
      // 左：铜钱三枚
      const left = document.createElement('div');
      left.className = 'trio-col coins-col';
      for (let i = 0; i < 3; i++) {
        const c = document.createElement('div');
        c.className = 'mini-coin';
        c.style.animationDelay = (i * 0.18) + 's';
        left.appendChild(c);
      }
      const lLabel = document.createElement('div');
      lLabel.className = 'trio-label';
      lLabel.textContent = '錢';
      left.appendChild(lLabel);

      // 中：月将转盘
      const mid = document.createElement('div');
      mid.className = 'trio-col wheel-col';
      const disc = document.createElement('div');
      disc.className = 'mini-wheel';
      ['子','寅','辰','午','申','戌'].forEach((z, i) => {
        const dot = document.createElement('div');
        dot.className = 'mini-wheel-dot';
        dot.style.transform = `rotate(${i * 60}deg) translateY(-32px)`;
        dot.textContent = z;
        disc.appendChild(dot);
      });
      const mc = document.createElement('div');
      mc.className = 'mini-wheel-center';
      mc.textContent = '將';
      disc.appendChild(mc);
      mid.appendChild(disc);
      const mLabel = document.createElement('div');
      mLabel.className = 'trio-label';
      mLabel.textContent = '將';
      mid.appendChild(mLabel);

      // 右：九宫
      const right = document.createElement('div');
      right.className = 'trio-col palaces-col';
      const grid = document.createElement('div');
      grid.className = 'mini-palaces';
      const stems = ['乙','丙','丁','戊','己','庚','辛','壬','癸'];
      stems.forEach((s, i) => {
        const c = document.createElement('div');
        c.className = 'mini-palace';
        if (i === 4) c.classList.add('center');
        c.style.animationDelay = (i * 0.06) + 's';
        c.textContent = s;
        grid.appendChild(c);
      });
      right.appendChild(grid);
      const rLabel = document.createElement('div');
      rLabel.className = 'trio-label';
      rLabel.textContent = '宮';
      right.appendChild(rLabel);

      container.appendChild(left);
      container.appendChild(mid);
      container.appendChild(right);
    }
  }

  async function stageMeditate(refs, question, opt) {
    const stage = refs.stageMeditate;
    refs.title1.textContent = L(opt.title1Key, opt.title1);
    refs.quote1.innerHTML = opt.quote1;
    refs.question.textContent = question ? '"' + question + '"' : '';
    stage.hidden = false;
    const words = ['三', '二', '一'];
    for (let i = 0; i < words.length; i++) {
      if (refs.skipped()) break;
      refs.countdown.textContent = words[i];
      refs.countdown.style.animation = 'none';
      void refs.countdown.offsetWidth;
      refs.countdown.style.animation = '';
      await sleep(900);
    }
    stage.hidden = true;
  }

  async function stageVisual(refs, fetchPromise, opt) {
    const stage = refs.stageVisual;
    refs.title2.textContent = L(opt.title2Key, opt.title2);
    refs.hint2.textContent = L(opt.hint2Key, opt.hint2);
    renderVisual(refs.visual, opt.visual);
    stage.hidden = false;
    const minMs = 2500;
    const [res] = await Promise.all([fetchPromise, sleep(minMs)]);
    stage.hidden = true;
    return res;
  }

  async function run({ kind, question, fetchPromise }) {
    const opt = PRESETS[kind];
    if (!opt) throw new Error('未知仪式类型: ' + kind);
    const o = ensureOverlay();
    const refs = {
      seal: o.querySelector('#ritualSeal'),
      stageMeditate: o.querySelector('#rStageMeditate'),
      stageVisual: o.querySelector('#rStageVisual'),
      title1: o.querySelector('#rTitle1'),
      quote1: o.querySelector('#rQuote1'),
      question: o.querySelector('#rQuestion'),
      countdown: o.querySelector('#rCountdown'),
      title2: o.querySelector('#rTitle2'),
      visual: o.querySelector('#rVisual'),
      hint2: o.querySelector('#rHint2'),
      skipBtn: o.querySelector('#rSkip'),
    };
    refs.seal.textContent = opt.seal;

    let _skipped = false;
    refs.skipped = () => _skipped;
    refs.skipBtn.onclick = () => { _skipped = true; };

    show(o);
    try {
      await stageMeditate(refs, question, opt);
      const res = await stageVisual(refs, fetchPromise, opt);
      hide(o);
      return res;
    } catch (e) {
      hide(o);
      throw e;
    }
  }

  window.Ritual = { run };
})();
