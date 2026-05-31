// AI 解卦前端模块：在「解卦提示词」下方挂一个「召 AI 解卦」按钮，
// 点击后把提示词 POST 到 /api/interpret，渲染模型解读。
//
// 用法：
//   window.AIInterpret.mount({
//     getPrompt: () => promptText,        // 返回当前提示词字符串
//     afterEl:   element or element id,   // 把解卦区插到这个元素之后
//     key:       'zhouyi',               // 唯一标识（避免重复挂载），可选
//   });
//
// /api/interpret 不消费访问码（起卦时已消费），故用普通 fetch；
// 服务端未配置 LLM 时返回 503，前端据此提示改用「复制提示词」。
(function () {
  'use strict';

  const tr = (k, f) => (window.I18n && window.I18n.t) ? window.I18n.t(k, f) : f;

  function escapeHtml(s) {
    if (s == null) return '';
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  // 行内格式：先转义 HTML，再处理 `代码`、**粗**、*斜*、~~删除~~、[链接](url)。
  function inlineMd(s) {
    let t = escapeHtml(s);
    t = t.replace(/`([^`]+)`/g, '<code>$1</code>');
    t = t.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    t = t.replace(/__([^_]+)__/g, '<strong>$1</strong>');
    t = t.replace(/~~([^~]+)~~/g, '<del>$1</del>');
    // 斜体：单 * 或单 _（避开已处理的 ** __）
    t = t.replace(/(^|[^*])\*([^*\n]+)\*(?!\*)/g, '$1<em>$2</em>');
    t = t.replace(/(^|[^_])_([^_\n]+)_(?!_)/g, '$1<em>$2</em>');
    // 链接（已转义，安全）：仅放行 http/https，避免 javascript: 等
    t = t.replace(/\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');
    return t;
  }

  // 轻量 Markdown 渲染器：支持标题、有序/无序列表、引用、分隔线、段落与行内格式。
  // 不引入第三方库；对流式中途的不完整标记（如半个 **）做容忍处理。
  function renderText(src) {
    const lines = String(src == null ? '' : src).replace(/\r\n?/g, '\n').split('\n');
    let html = '';
    let listType = null; // 'ul' | 'ol'
    let para = [];
    const closeList = () => { if (listType) { html += '</' + listType + '>'; listType = null; } };
    const flushPara = () => {
      if (para.length) { html += '<p>' + para.map(inlineMd).join('<br>') + '</p>'; para = []; }
    };
    // 表格分隔行：如 |---|:--:|---| —— 仅由 | : - 和空格组成，且含 - 与 |
    const isTableSep = (s) => {
      const x = s.trim();
      return x.indexOf('|') !== -1 && x.indexOf('-') !== -1 && /^[|:\-\s]+$/.test(x);
    };
    // 拆分表格行为单元格（去掉首尾的 |）
    const splitRow = (s) => {
      let x = s.trim();
      if (x.startsWith('|')) x = x.slice(1);
      if (x.endsWith('|')) x = x.slice(0, -1);
      return x.split('|').map((c) => c.trim());
    };

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const t = line.trim();
      let m;

      if (/^(-{3,}|\*{3,}|_{3,})$/.test(t)) {            // 分隔线
        flushPara(); closeList(); html += '<hr>'; continue;
      }
      if ((m = /^(#{1,6})\s+(.*)$/.exec(t))) {            // 标题（卡片内压成 h4/h5）
        flushPara(); closeList();
        const tag = m[1].length <= 2 ? 'h4' : 'h5';
        html += '<' + tag + '>' + inlineMd(m[2]) + '</' + tag + '>'; continue;
      }
      if ((m = /^>\s?(.*)$/.exec(t))) {                   // 引用
        flushPara(); closeList();
        html += '<blockquote>' + inlineMd(m[1]) + '</blockquote>'; continue;
      }
      // 表格：当前行含 | 且下一行是分隔行（|---|---|）
      if (t.indexOf('|') !== -1 && i + 1 < lines.length && isTableSep(lines[i + 1])) {
        flushPara(); closeList();
        const headers = splitRow(t);
        const rows = [];
        let j = i + 2;
        while (j < lines.length && lines[j].trim() !== '' && lines[j].indexOf('|') !== -1) {
          rows.push(splitRow(lines[j]));
          j++;
        }
        html += '<table class="md-table"><thead><tr>' +
          headers.map((h) => '<th>' + inlineMd(h) + '</th>').join('') +
          '</tr></thead><tbody>' +
          rows.map((r) => '<tr>' +
            headers.map((_, k) => '<td>' + inlineMd(r[k] || '') + '</td>').join('') +
            '</tr>').join('') +
          '</tbody></table>';
        i = j - 1; // 循环末尾会 i++，使其指向 j
        continue;
      }
      if ((m = /^\d+[.)]\s+(.*)$/.exec(t))) {             // 有序列表
        flushPara();
        if (listType !== 'ol') { closeList(); html += '<ol>'; listType = 'ol'; }
        html += '<li>' + inlineMd(m[1]) + '</li>'; continue;
      }
      if ((m = /^[-*+]\s+(.*)$/.exec(t))) {               // 无序列表
        flushPara();
        if (listType !== 'ul') { closeList(); html += '<ul>'; listType = 'ul'; }
        html += '<li>' + inlineMd(m[1]) + '</li>'; continue;
      }
      if (t === '') { flushPara(); closeList(); continue; } // 空行 → 分段

      closeList();                                         // 普通文本行
      para.push(line);
    }
    flushPara(); closeList();
    return html;
  }

  function resolveEl(elOrId) {
    if (!elOrId) return null;
    return (typeof elOrId === 'string') ? document.getElementById(elOrId) : elOrId;
  }

  // readSSE 读取一个 ReadableStream（SSE 流），逐个事件回调 onEvent(event, dataStr)。
  // 解析以空行分隔的事件块，取其中的 event: 与 data: 字段（data 多行则拼接）。
  async function readSSE(stream, onEvent) {
    const reader = stream.getReader();
    const decoder = new TextDecoder('utf-8');
    let buf = '';
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      buf += decoder.decode(value, { stream: true });
      // 按事件块（\n\n）切分，保留最后一个可能不完整的块
      let idx;
      while ((idx = buf.indexOf('\n\n')) !== -1) {
        const block = buf.slice(0, idx);
        buf = buf.slice(idx + 2);
        dispatchBlock(block, onEvent);
      }
    }
    if (buf.trim()) dispatchBlock(buf, onEvent);
  }

  function dispatchBlock(block, onEvent) {
    let event = 'message';
    const dataLines = [];
    block.split('\n').forEach((line) => {
      if (line.startsWith('event:')) event = line.slice(6).trim();
      else if (line.startsWith('data:')) dataLines.push(line.slice(5).replace(/^ /, ''));
    });
    if (dataLines.length) onEvent(event, dataLines.join('\n'));
  }

  function mount(opts) {
    opts = opts || {};
    const after = resolveEl(opts.afterEl);
    if (!after || typeof opts.getPrompt !== 'function') return;

    const key = opts.key || 'default';
    const blockId = 'aiInterpretBlock-' + key;

    // 已挂载过则只复用（保留上次解读，不重复创建）
    let block = document.getElementById(blockId);
    if (!block) {
      block = document.createElement('div');
      block.id = blockId;
      block.className = 'ai-interpret-block';
      block.innerHTML = `
        <div class="prompt-toolbar" style="margin-top:14px;">
          <button type="button" class="secondary ai-interpret-btn">${tr('interpret.btn', '召 AI 解卦')}</button>
          <span class="copy-status ai-interpret-status"></span>
        </div>
        <div class="ai-interpret-result hex-info" hidden style="margin-top:12px;line-height:1.9;"></div>
      `;
      after.parentNode.insertBefore(block, after.nextSibling);

      const btn = block.querySelector('.ai-interpret-btn');
      const status = block.querySelector('.ai-interpret-status');
      const resultBox = block.querySelector('.ai-interpret-result');

      btn.addEventListener('click', async () => {
        const prompt = (opts.getPrompt() || '').trim();
        if (!prompt) {
          status.textContent = tr('interpret.noprompt', '尚無提示詞，請先起卦');
          status.classList.add('show');
          setTimeout(() => status.classList.remove('show'), 1800);
          return;
        }

        btn.disabled = true;
        const oldLabel = btn.textContent;
        btn.textContent = tr('interpret.loading', 'AI 正在解卦……');
        resultBox.hidden = false;
        resultBox.innerHTML = `<p style="color:var(--ink-soft);">${tr('interpret.loading.hint', '正請大模型研讀卦象，逐字顯現……')}</p>`;

        // —— 打字机逐字显示 ——
        // acc：已收到的全文；shown：已显示字符数；streamDone：SSE 是否结束。
        // 定时器每 ~18ms 把 step 个字从 acc 吐到屏幕；落后越多吐越快，避免长文拖尾。
        let acc = '';
        let shown = 0;
        let started = false;
        let streamDone = false;
        let timer = null;
        const caret = '<span class="ai-caret"></span>';

        const paint = () => {
          const finished = streamDone && shown >= acc.length;
          resultBox.innerHTML = renderText(acc.slice(0, shown)) + (finished ? '' : caret);
        };
        const ensureTyper = () => {
          if (timer) return;
          timer = setInterval(() => {
            if (shown < acc.length) {
              const remain = acc.length - shown;
              const step = remain > 120 ? 8 : (remain > 40 ? 4 : (remain > 12 ? 2 : 1));
              shown = Math.min(acc.length, shown + step);
              paint();
            } else if (streamDone) {
              clearInterval(timer); timer = null;
              paint();
            }
          }, 18);
        };
        // 等打字机把已收到的内容吐完
        const drain = () => new Promise((resolve) => {
          const check = () => (shown >= acc.length ? resolve() : setTimeout(check, 24));
          check();
        });

        try {
          const res = await fetch('/api/interpret', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ prompt }),
          });

          // 未配置 / 错误：服务端在开流前以 JSON 返回（非 event-stream）
          const ctype = res.headers.get('Content-Type') || '';
          if (!ctype.includes('text/event-stream')) {
            let data = {};
            try { data = await res.json(); } catch (_) {}
            if (res.status === 503) {
              resultBox.innerHTML = `<p style="color:var(--vermilion-dark);">${tr('interpret.unconfigured', '服務端未配置解卦 API。可點上方「複製提示詞」，自行貼到任意 AI 解讀。')}</p>`;
            } else {
              resultBox.innerHTML = `<p style="color:var(--vermilion-dark);">${tr('interpret.failed', '解卦失敗')}：${escapeHtml(data.error || res.statusText)}</p>`;
            }
            return;
          }

          // 读取 SSE 流：每段只追加到 acc，由打字机定时器逐字吐出
          let streamErr = '';
          await readSSE(res.body, (event, dataStr) => {
            let payload = {};
            try { payload = JSON.parse(dataStr); } catch (_) {}
            if (event === 'delta' && payload.text) {
              if (!started) { started = true; resultBox.innerHTML = ''; }
              acc += payload.text;
              ensureTyper();
            } else if (event === 'error') {
              streamErr = payload.error || tr('interpret.failed', '解卦失敗');
            }
          });

          streamDone = true;
          await drain();                       // 等逐字吐完
          if (timer) { clearInterval(timer); timer = null; }
          paint();                             // 去掉光标的最终渲染

          if (streamErr) {
            resultBox.innerHTML = renderText(acc) +
              `<p style="color:var(--vermilion-dark);margin-top:10px;">${escapeHtml(streamErr)}</p>`;
          } else if (!acc.trim()) {
            resultBox.innerHTML = `<p style="color:var(--vermilion-dark);">${tr('interpret.empty', 'AI 未返回內容，請重試')}</p>`;
          }
        } catch (e) {
          streamDone = true;
          if (timer) { clearInterval(timer); timer = null; }
          const tail = acc.trim() ? renderText(acc) + '<hr style="opacity:.3;margin:10px 0;">' : '';
          resultBox.innerHTML = tail + `<p style="color:var(--vermilion-dark);">${tr('interpret.neterr', '網絡錯誤')}：${escapeHtml(e.message || String(e))}</p>`;
        } finally {
          if (timer) { clearInterval(timer); timer = null; }
          btn.disabled = false;
          btn.textContent = oldLabel;
        }
      });
    }

    // 每次起新卦后清掉旧解读，避免张冠李戴
    const resultBox = block.querySelector('.ai-interpret-result');
    if (resultBox) { resultBox.hidden = true; resultBox.innerHTML = ''; }
  }

  window.AIInterpret = { mount };
})();
