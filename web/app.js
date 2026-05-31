// 周易筮占前端 —— 含起卦仪式
(function () {
  'use strict';

  const $ = (id) => document.getElementById(id);
  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

  const state = { method: 'yarrow', busy: false, skip: false };

  // ---------- 初始化 ----------
  document.addEventListener('DOMContentLoaded', async () => {
    bindTabs();
    bindDivine();
    bindCopy();
    bindRitualSkip();
    await loadQuestionTypes();
  });

  function bindTabs() {
    document.querySelectorAll('.method-tabs .tab').forEach((btn) => {
      btn.addEventListener('click', () => {
        document.querySelectorAll('.method-tabs .tab').forEach((b) => b.classList.remove('active'));
        btn.classList.add('active');
        state.method = btn.dataset.method;
        $('numberFields').hidden = state.method !== 'number';
      });
    });
  }

  function bindDivine() {
    $('divineBtn').addEventListener('click', () => {
      if (!state.busy) doDivine();
    });
  }

  function bindRitualSkip() {
    $('ritualSkip').addEventListener('click', () => { state.skip = true; });
  }

  function bindCopy() {
    const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;
    $('copyPromptBtn').addEventListener('click', () => {
      const t = $('promptText');
      t.select();
      t.setSelectionRange(0, t.value.length);
      let ok = false;
      try { ok = document.execCommand && document.execCommand('copy'); } catch (_) {}
      if (!ok && navigator.clipboard) {
        navigator.clipboard.writeText(t.value).then(() => flashCopy(tr('msg.copied', '已複製到剪貼板')));
      } else {
        flashCopy(ok ? tr('msg.copied', '已複製到剪貼板') : tr('msg.copy.failed', '複製失敗，請手動選擇'));
      }
      window.getSelection().removeAllRanges();
    });
  }
  function flashCopy(msg) {
    const s = $('copyStatus');
    s.textContent = msg;
    s.classList.add('show');
    setTimeout(() => s.classList.remove('show'), 1800);
  }

  async function loadQuestionTypes() {
    try {
      const res = await fetch('/api/question-types');
      const list = await res.json();
      const sel = $('questionType');
      sel.innerHTML = '';
      const placeholder = document.createElement('option');
      placeholder.value = '';
      placeholder.textContent = '— 不指定 —';
      sel.appendChild(placeholder);
      list.forEach((t) => {
        const opt = document.createElement('option');
        opt.value = t.value;
        opt.textContent = t.label;
        sel.appendChild(opt);
      });
    } catch (e) {
      const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;
      console.error(tr('msg.qtype.load.fail', '加載問題類型失敗'), e);
    }
  }

  // ---------- 起卦主流程 ----------
  async function doDivine() {
    const btn = $('divineBtn');
    btn.disabled = true;
    state.busy = true;
    state.skip = false;

    const body = {
      method: state.method,
      question: $('question').value.trim(),
      questionType: $('questionType').value,
    };
    if (state.method === 'number') {
      body.upper = parseInt($('upper').value, 10);
      body.lower = parseInt($('lower').value, 10);
      body.changing = parseInt($('changing').value, 10);
      if (!Number.isFinite(body.upper) || !Number.isFinite(body.lower)) {
        const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;
        showStatus(tr('msg.fill.upper.lower', '請填寫上卦數與下卦數（整數即可）'), true);
        btn.disabled = false;
        state.busy = false;
        return;
      }
      if (!Number.isFinite(body.changing)) body.changing = 0;
    }

    $('divineStatus').hidden = true;

    // 开启仪式
    showOverlay();
    try {
      // 阶段一：凝神（约 3 秒）+ 同时并发请求
      body.longitude = await Location.get();
      const fetchPromise = AccessCode.authedFetch('/api/divine', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      await stageMeditate(body.question);

      // 阶段二：摇卦动画（至少 2.5 秒），直到请求返回
      const res = await stageShaking(fetchPromise);
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        hideOverlay();
        const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;
        showStatus(tr('msg.cast.failed', '起卦失敗') + '：' + (err.error || res.statusText), true);
        return;
      }
      const data = await res.json();

      // 阶段三：逐爻揭示（约 3.2 秒）
      await stageReveal(data);

      hideOverlay();

      // 渲染最终结果
      renderResult(data);
    } catch (e) {
      hideOverlay();
      const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;
      showStatus(tr('msg.network.error', '網絡錯誤') + '：' + e.message, true);
    } finally {
      btn.disabled = false;
      state.busy = false;
    }
  }

  function showStatus(msg, isError) {
    const status = $('divineStatus');
    status.hidden = false;
    status.textContent = msg;
    status.style.borderLeftColor = isError ? 'var(--vermilion)' : '';
  }

  // ---------- 仪式：遮罩与阶段 ----------
  function showOverlay() {
    const o = $('ritualOverlay');
    o.hidden = false;
    document.body.style.overflow = 'hidden';
    ['stageMeditate', 'stageShaking', 'stageReveal'].forEach((id) => { $(id).hidden = true; });
  }
  function hideOverlay() {
    $('ritualOverlay').hidden = true;
    document.body.style.overflow = '';
  }

  async function stageMeditate(question) {
    const stage = $('stageMeditate');
    const qEl = $('ritualQuestion');
    const cd = $('ritualCountdown');
    qEl.textContent = question ? '“' + question + '”' : '';
    stage.hidden = false;

    const words = ['三', '二', '一'];
    for (let i = 0; i < words.length; i++) {
      if (state.skip) break;
      cd.textContent = words[i];
      cd.style.animation = 'none';
      void cd.offsetWidth;
      cd.style.animation = '';
      await sleep(900);
    }
    stage.hidden = true;
  }

  async function stageShaking(fetchPromise) {
    const stage = $('stageShaking');
    const title = $('shakingTitle');
    const hint = $('shakingHint');
    const coinStage = $('coinStage');

    // 根据方法换文案
    const map = {
      coin:   { title: '搖　卦', hint: '三錢入掌，六擲成卦……', visual: 'coins' },
      yarrow: { title: '揲　蓍', hint: '五十蓍草，十八變而成卦……', visual: 'yarrow' },
      number: { title: '布　卦', hint: '以數定象，上下相合……', visual: 'numbers' },
    };
    const cfg = map[state.method] || map.coin;
    title.textContent = cfg.title;
    hint.textContent = cfg.hint;

    renderShakingVisual(coinStage, cfg.visual);

    stage.hidden = false;
    const minMs = 2500;
    const [res] = await Promise.all([fetchPromise, sleep(minMs)]);
    stage.hidden = true;
    return res;
  }

  function renderShakingVisual(container, kind) {
    container.innerHTML = '';
    container.className = 'ritual-coins';
    if (kind === 'coins') {
      for (let i = 0; i < 3; i++) {
        const c = document.createElement('div');
        c.className = 'coin';
        container.appendChild(c);
      }
    } else if (kind === 'yarrow') {
      container.classList.add('yarrow');
      for (let i = 0; i < 7; i++) {
        const s = document.createElement('div');
        s.className = 'stalk';
        s.style.animationDelay = (i * 0.12) + 's';
        container.appendChild(s);
      }
    } else if (kind === 'numbers') {
      container.classList.add('numbers');
      ['乾', '兌', '離', '震', '巽', '坎', '艮', '坤'].forEach((t, i) => {
        const n = document.createElement('div');
        n.className = 'numglyph';
        n.textContent = t;
        n.style.animationDelay = (i * 0.08) + 's';
        container.appendChild(n);
      });
    }
  }

  async function stageReveal(data) {
    const stage = $('stageReveal');
    const wrap = $('revealHex');
    const hint = $('revealHint');
    wrap.innerHTML = '';
    const posNames = ['初', '二', '三', '四', '五', '上'];

    // 从上到下布局（视觉上上爻在顶），但按"从初爻到上爻"的顺序依次点亮
    // 先把 6 行占位按视觉顺序渲染
    const rows = [];
    for (let i = 5; i >= 0; i--) {
      const ln = data.lines[i];
      const row = document.createElement('div');
      row.className = 'row';
      const bar = document.createElement('div');
      bar.className = 'bar ' + (ln.isYang ? 'yang' : 'yin') + (ln.isChange ? ' change' : '');
      const lbl = document.createElement('div');
      lbl.className = 'label';
      lbl.textContent = posNames[i] + (ln.isYang ? '九' : '六') + (ln.isChange ? ' ★' : '');
      row.appendChild(bar);
      row.appendChild(lbl);
      wrap.appendChild(row);
      rows.push(row);
    }
    // rows[0]=上爻视觉位置，rows[5]=初爻视觉位置
    stage.hidden = false;

    hint.textContent = '自初爻至上爻，次第顯現……';

    for (let i = 0; i < 6; i++) {
      if (state.skip) {
        rows.forEach((r) => r.classList.add('reveal'));
        break;
      }
      // 初爻=i=0 对应 rows[5]
      rows[5 - i].classList.add('reveal');
      await sleep(480);
    }

    // 最后显示卦名提示
    const h = data.mainHex;
    if (h) {
      hint.innerHTML = `卦　成　——　<span style="color:#e8c97a;font-family:'STKaiti',serif;font-size:1.2rem;letter-spacing:.2em;">${h.symbol} ${h.name}卦</span>（第${h.number}卦）`;
    }
    await sleep(state.skip ? 200 : 900);
    stage.hidden = true;
  }

  // ---------- 渲染结果 ----------
  function renderResult(data) {
    $('resultSection').hidden = false;
    $('promptSection').hidden = false;

    renderHexArt(data);
    renderHexMeta(data);
    renderLines(data);
    renderHexInfo('mainHexInfo', data.mainHex);

    if (data.changeHex) {
      $('changeHexBlock').hidden = false;
      renderHexInfo('changeHexInfo', data.changeHex);
    } else {
      $('changeHexBlock').hidden = true;
    }

    renderDerived(data.derived);
    renderTiming(data.timing);

    $('guide').textContent = data.guide || '';
    $('promptText').value = data.prompt || '';

    // AI 解卦按钮（挂在提示词文本框之后）
    if (window.AIInterpret) {
      window.AIInterpret.mount({
        key: 'zhouyi',
        afterEl: 'promptText',
        getPrompt: () => $('promptText').value,
      });
    }

    // 印心小卡（结果区顶部）
    if (window.Journal && window.Journal.renderYinXin) {
      const head = (data.mainHex ? data.mainHex.name + '卦' : '') +
        (data.changeHex ? ' → ' + data.changeHex.name + '卦' : '');
      window.Journal.renderYinXin({
        kind: 'zhouyi',
        data,
        mountId: 'yinxinSlot',
        question: data.question || $('question').value.trim(),
        questionType: $('questionType').value,
        headline: head,
        ganzhiTime: (data.timing && data.timing.ganZhiTime) || ''
      });
    }

    setTimeout(() => {
      $('resultSection').scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 100);
  }

  function renderHexArt(data) {
    const el = $('hexArt');
    el.innerHTML = '';
    const posNames = ['初', '二', '三', '四', '五', '上'];
    for (let i = 5; i >= 0; i--) {
      const ln = data.lines[i];
      const row = document.createElement('div');
      row.className = 'row';
      const bar = document.createElement('div');
      bar.className = 'bar ' + (ln.isYang ? 'yang' : 'yin') + (ln.isChange ? ' change' : '');
      if (ln.isChange) {
        const mk = document.createElement('span');
        mk.className = 'mark';
        mk.textContent = ln.isYang ? '○' : '×';
        bar.appendChild(mk);
      }
      const label = document.createElement('span');
      const yy = ln.isYang ? '九' : '六';
      label.textContent = posNames[i] + yy + (ln.isChange ? '　變' : '');
      row.appendChild(bar);
      row.appendChild(label);
      el.appendChild(row);
    }
    // 交错进入
    setTimeout(() => {
      el.querySelectorAll('.row').forEach((r, idx) => {
        setTimeout(() => r.classList.add('reveal'), idx * 90);
      });
    }, 50);
  }

  function renderHexMeta(data) {
    const h = data.mainHex;
    if (!h) { $('hexMeta').innerHTML = ''; return; }
    $('hexMeta').innerHTML = `
      <div class="hex-name"><span class="hex-symbol">${h.symbol}</span>${h.name}</div>
      <div class="hex-sub">第 ${h.number} 卦 · ${escapeHtml(data.methodLabel)}</div>
      <div class="hex-sub small">上卦 ${h.upper.symbol} ${h.upper.name}（${h.upper.nature}） / 下卦 ${h.lower.symbol} ${h.lower.name}（${h.lower.nature}）</div>
      ${data.question ? `<div class="hex-sub small">所問：${escapeHtml(data.question)}</div>` : ''}
      ${data.questionLabel ? `<div class="hex-sub small">問類：${escapeHtml(data.questionLabel)}</div>` : ''}
    `;
  }

  function renderLines(data) {
    const body = $('linesBody');
    body.innerHTML = '';
    const posNames = ['初', '二', '三', '四', '五', '上'];
    for (let i = 5; i >= 0; i--) {
      const ln = data.lines[i];
      const tr = document.createElement('tr');
      if (ln.isChange) tr.classList.add('changing');
      tr.innerHTML = `
        <td>${posNames[i]}${ln.isYang ? '九' : '六'}${ln.isChange ? ' ★' : ''}</td>
        <td>${ln.value}</td>
        <td>${ln.typeName}</td>
        <td class="sym">${ln.symbol}</td>
        <td>${escapeHtml(ln.text || '')}</td>
      `;
      body.appendChild(tr);
    }
  }

  function renderHexInfo(id, h) {
    if (!h) { $(id).innerHTML = ''; return; }
    const linesHtml = h.lines
      .slice()
      .reverse()
      .map((l) => `<p><span class="label">${l.type}</span>${escapeHtml(l.text)}</p>`)
      .join('');
    $(id).innerHTML = `
      <div class="hex-info">
        <p><span class="label">卦名</span>${h.symbol} ${h.name}（第 ${h.number} 卦）</p>
        <p><span class="label">上下卦</span>${h.upper.name}（${h.upper.nature}） / ${h.lower.name}（${h.lower.nature}）</p>
        <p><span class="label">卦辭</span>${escapeHtml(h.judgment)}</p>
        <p><span class="label">象辭</span>${escapeHtml(h.image)}</p>
        <h4>爻辭</h4>
        ${linesHtml}
      </div>
    `;
  }

  function renderDerived(d) {
    const grid = $('derivedGrid');
    grid.innerHTML = '';
    const items = [
      { kind: '互卦', desc: '内部之機', h: d.mutual },
      { kind: '錯卦', desc: '相對之面', h: d.opposite },
      { kind: '綜卦', desc: '換位之觀', h: d.reverse },
    ];
    items.forEach((it) => {
      if (!it.h) return;
      const card = document.createElement('div');
      card.className = 'derived-card';
      card.innerHTML = `
        <div class="kind">${it.kind} · ${it.desc}</div>
        <div class="name"><span class="sym">${it.h.symbol}</span>${it.h.name}（第 ${it.h.number} 卦）</div>
        <div class="judg">${escapeHtml(it.h.judgment)}</div>
      `;
      grid.appendChild(card);
    });
  }

  function renderTiming(ti) {
    if (!ti) { $('timingBlock').hidden = true; return; }
    $('timingBlock').hidden = false;
    const parts = [];
    parts.push(`<p><span class="label">陽曆</span>${escapeHtml(ti.solarTime)}</p>`);
    parts.push(`<p><span class="label">陰曆</span>${escapeHtml(ti.lunarDesc)}</p>`);
    parts.push(`<p><span class="label">干支</span>${escapeHtml(ti.ganZhiSummary)}</p>`);
    if (ti.jieQiName) {
      const next = ti.nextJieQiName ? `，距「${ti.nextJieQiName}」（${ti.nextJieQiDate}）` : '';
      parts.push(`<p><span class="label">節氣</span>${escapeHtml(ti.jieQiDay)}${next}</p>`);
    }
    if (ti.monthlyHex) {
      parts.push(`<p><span class="label">時令卦</span>第 ${ti.monthlyHex.number} 卦 ${ti.monthlyHex.name} ${ti.monthlyHex.symbol} —— ${escapeHtml(ti.monthlyHexNote)}</p>`);
    }
    $('timingInfo').innerHTML = `<div class="hex-info">${parts.join('')}</div>`;
  }

  function escapeHtml(s) {
    if (s == null) return '';
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }
})();
