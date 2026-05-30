// i18n: 繁體 / 简体 / English 切换
//
// 原则：
//   - 源 HTML 一律繁體中文（zh-TW/zh-Hant），這是母版
//   - zh-CN 模式：在客戶端把可見文本通過字典做繁→簡字符級轉換
//   - en 模式：data-i18n="key" 元素查 EN 字典；data-i18n-keep 元素保留中文
//   - 動態插入的 DOM（app.js / ritual.js / 各頁內聯 JS 渲染的盤面結果）
//     由 MutationObserver 自動處理；同時暴露 I18n.apply(node) 供主動觸發
//
// 用法：
//   <script src="/i18n.js"></script>  // 必須在 access.js / app.js 之前或之後均可
//   切換按鈕：<a data-lang-switch="zh-TW">繁</a> 等

(function () {
  'use strict';

  const STORAGE_KEY = 'zhouyi.lang';
  const DEFAULT_LANG = 'zh-TW';
  const SUPPORTED = ['zh-TW', 'zh-CN', 'zh-mod'];

  // 當前語言。先讀 localStorage，否則默認繁體。
  function readLang() {
    try {
      const v = localStorage.getItem(STORAGE_KEY);
      if (v && SUPPORTED.indexOf(v) >= 0) return v;
    } catch (_) {}
    return DEFAULT_LANG;
  }
  let curLang = readLang();

  // 對外暴露 lang，便於其他模塊判斷當前語言
  function getLang() { return curLang; }

  // ===== 繁→簡字符映射表 =====
  // 覆蓋本系統用到的字（卦辭、爻辭、神煞、占斷常見字）。逐字一對一，運行時 O(N) 處理。
  // 注意：少數字繁→簡有歧義（如「乾/干」「後/后」），這裡採取「項目語境下最常見義」的取捨。
  // 來源：基於常用 OpenCC TWVariants → Simplified 子集 + 本項目術語人工校對。
  const T2S = {
    // 由 OpenCC t2s 生成的繁→简字符級映射；專名例外（乾、後）保留繁體
    '並':'并','亂':'乱','併':'并','來':'来','係':'系','個':'个','倫':'伦','備':'备',
    '傳':'传','傾':'倾','僅':'仅','儀':'仪','償':'偿','優':'优','兌':'兑','內':'内',
    '兩':'两','別':'别','刪':'删','則':'则','創':'创','劃':'划','動':'动','務':'务',
    '勝':'胜','勞':'劳','勢':'势','匯':'汇','卻':'却','參':'参','員':'员','問':'问',
    '啟':'启','單':'单','嗎':'吗','嚴':'严','囂':'嚣','圍':'围','執':'执','場':'场',
    '學':'学','宮':'宫','實':'实','審':'审','寫':'写','寬':'宽','將':'将','專':'专',
    '尋':'寻','對':'对','導':'导','屬':'属','師':'师','幫':'帮','幾':'几','庫':'库',
    '張':'张','強':'强','後':'后','從':'从','復':'复','恆':'恒','惡':'恶','惱':'恼',
    '愛':'爱','態':'态','應':'应','懸':'悬','戰':'战','戲':'戏','戶':'户','捨':'舍',
    '採':'采','換':'换','損':'损','搖':'摇','撥':'拨','擇':'择','擊':'击','擔':'担',
    '據':'据','擲':'掷','擾':'扰','敗':'败','數':'数','斷':'断','於':'于','時':'时',
    '晝':'昼','暫':'暂','曆':'历','書':'书','會':'会','東':'东','條':'条','業':'业',
    '極':'极','構':'构','標':'标','樣':'样','機':'机','檢':'检','檻':'槛','權':'权',
    '歡':'欢','歷':'历','殺':'杀','氣':'气','決':'决','沒':'没','況':'况','淨':'净',
    '淺':'浅','測':'测','準':'准','滯':'滞','瀆':'渎','瀏':'浏','為':'为','無':'无',
    '煙':'烟','狀':'状','獨':'独','獲':'获','現':'现','瑣':'琐','產':'产','畢':'毕',
    '異':'异','當':'当','療':'疗','發':'发','盡':'尽','盤':'盘','盪':'荡','眾':'众',
    '確':'确','碼':'码','稱':'称','積':'积','穩':'稳','窺':'窥','節':'节','範':'范',
    '篩':'筛','簡':'简','簽':'签','籤':'签','約':'约','紐':'纽','細':'细','終':'终',
    '結':'结','絡':'络','給':'给','統':'统','經':'经','綜':'综','維':'维','網':'网',
    '緒':'绪','緣':'缘','編':'编','緩':'缓','緻':'致','縱':'纵','總':'总','織':'织',
    '繫':'系','繼':'继','續':'续','義':'义','習':'习','聖':'圣','聲':'声','職':'职',
    '膩':'腻','臨':'临','與':'与','萬':'万','蓋':'盖','藥':'药','處':'处','號':'号',
    '術':'术','補':'补','裡':'里','製':'制','複':'复','褻':'亵','見':'见','規':'规',
    '視':'视','親':'亲','覽':'览','觀':'观','計':'计','訊':'讯','託':'托','記':'记',
    '訟':'讼','訪':'访','設':'设','許':'许','訴':'诉','診':'诊','註':'注','詐':'诈',
    '評':'评','詞':'词','詢':'询','試':'试','話':'话','該':'该','詳':'详','認':'认',
    '語':'语','誠':'诚','誤':'误','誰':'谁','課':'课','談':'谈','請':'请','論':'论',
    '諮':'咨','諸':'诸','謀':'谋','講':'讲','謹':'谨','證':'证','譏':'讥','識':'识',
    '議':'议','讀':'读','變':'变','讓':'让','豐':'丰','豬':'猪','負':'负','財':'财',
    '貪':'贪','責':'责','貳':'贰','貴':'贵','買':'买','貼':'贴','資':'资','賠':'赔',
    '賢':'贤','賦':'赋','賭':'赌','賴':'赖','趨':'趋','軌':'轨','較':'较','載':'载',
    '輒':'辄','輔':'辅','輕':'轻','輪':'轮','輸':'输','轉':'转','辦':'办','辭':'辞',
    '這':'这','連':'连','進':'进','遊':'游','運':'运','過':'过','違':'违','遠':'远',
    '適':'适','遷':'迁','選':'选','還':'还','邊':'边','醫':'医','釋':'释','釣':'钓',
    '鈕':'钮','銅':'铜','錄':'录','錢':'钱','錯':'错','鍵':'键','鏡':'镜','鑑':'鉴',
    '鑒':'鉴','長':'长','門':'门','閉':'闭','開':'开','閏':'闰','閒':'闲','間':'间',
    '閱':'阅','闊':'阔','關':'关','闡':'阐','陰':'阴','陳':'陈','陽':'阳','際':'际',
    '隨':'随','險':'险','隱':'隐','雞':'鸡','離':'离','難':'难','電':'电','靈':'灵',
    '靜':'静','響':'响','頁':'页','項':'项','順':'顺','須':'须','預':'预','頭':'头',
    '題':'题','願':'愿','類':'类','顯':'显','風':'风','飛':'飞','飢':'饥','飯':'饭',
    '飽':'饱','養':'养','餵':'喂','馬':'马','驕':'骄','驗':'验','驚':'惊','驛':'驿',
    '體':'体','鬆':'松','鳥':'鸟','麼':'么','黃':'黄','點':'点','齊':'齐','龍':'龙'
  };

  // 一遍掃描即可（O(N)）。對未在字典中的字符保留原樣。
  function t2s(s) {
    if (!s) return s;
    let out = '';
    for (let i = 0; i < s.length; i++) {
      const ch = s[i];
      out += T2S[ch] || ch;
    }
    return out;
  }

  // 公開以便其他模塊（access.js / location.js / 內聯 JS）按當前語言生成消息時調用
  function tr(zhTW) {
    if (curLang === 'zh-CN') return t2s(zhTW);
    if (curLang === 'zh-mod') {
      // 「現代漢語」模式下，動態消息沒有 i18n key 就走 t2s 兜底（變簡體）
      return t2s(zhTW);
    }
    return zhTW; // zh-TW 原樣返回
  }

  // ===== EN 字典：UI 殼文本 =====
  // 由 i18n-en.js 注入到 window.__I18N_EN__；此處先聲明 fallback。
  function lookupEN(key) {
    const dict = window.__I18N_EN__ || {};
    return dict[key];
  }

  // ===== DOM 處理 =====
  // 1) 收集所有 data-i18n 元素 → key 替換（EN 模式）
  // 2) 對其餘文本節點：EN 模式下若處於 data-i18n-keep 子樹則保留中文，否則嘗試走 EN 標題（只應用於有 key 的元素）
  // 3) zh-CN 模式：對所有可見文本節點做 t2s

  function isInsideKeep(node) {
    let p = node.nodeType === 1 ? node : node.parentElement;
    while (p) {
      if (p.hasAttribute && p.hasAttribute('data-i18n-keep')) return true;
      p = p.parentElement;
    }
    return false;
  }

  // 文本節點原文緩存：第一次見到時記錄繁體原文，後續切換語言時始終以原文為基準
  const ORIG = new WeakMap();
  function origText(node) {
    if (ORIG.has(node)) return ORIG.get(node);
    ORIG.set(node, node.nodeValue);
    return node.nodeValue;
  }

  // EN 模式下：對某些屬性（placeholder/title/aria-label/value）也要處理
  const ATTRS = ['placeholder', 'title', 'aria-label'];
  const ORIG_ATTR = new WeakMap();
  function origAttr(el, name) {
    let m = ORIG_ATTR.get(el);
    if (!m) { m = {}; ORIG_ATTR.set(el, m); }
    if (!(name in m)) m[name] = el.getAttribute(name);
    return m[name];
  }

  function applyToNode(root) {
    if (!root) return;

    // ① 處理帶 data-i18n 的元素（替換 textContent）
    const i18nNodes = (root.nodeType === 1 && root.hasAttribute && root.hasAttribute('data-i18n'))
      ? [root]
      : (root.querySelectorAll ? Array.from(root.querySelectorAll('[data-i18n]')) : []);
    i18nNodes.forEach(el => {
      const key = el.getAttribute('data-i18n');
      // 對於用 data-i18n-html 標記的元素，原文用 innerHTML 緩存，譯文以 HTML 注入；
      // 否則用 textContent（純文本）
      const useHTML = el.hasAttribute('data-i18n-html');
      const origRead = useHTML ? el.innerHTML : el.textContent;
      const origUsed = el.getAttribute('data-i18n-orig') || origRead;
      if (!el.hasAttribute('data-i18n-orig')) el.setAttribute('data-i18n-orig', origUsed);
      let val = origUsed;
      if (curLang === 'zh-mod') {
        const mod = lookupEN(key);
        if (mod) val = mod;
        else val = t2s(origUsed); // 缺譯回退到簡體
      } else if (curLang === 'zh-CN') {
        val = t2s(origUsed);
      }
      if (useHTML) el.innerHTML = val;
      else el.textContent = val;
    });

    // ② 處理屬性（placeholder / title / aria-label）
    const attrNodes = root.querySelectorAll
      ? Array.from(root.querySelectorAll('[data-i18n-attr]'))
      : [];
    if (root.nodeType === 1 && root.hasAttribute && root.hasAttribute('data-i18n-attr')) attrNodes.unshift(root);
    attrNodes.forEach(el => {
      const spec = el.getAttribute('data-i18n-attr'); // 形如 "placeholder:zhouyi.placeholder"
      spec.split(',').forEach(pair => {
        const [attr, key] = pair.split(':').map(s => s.trim());
        const orig = origAttr(el, attr);
        if (orig == null) return;
        let val = orig;
        if (curLang === 'zh-mod') {
          const mod = lookupEN(key);
          val = mod || t2s(orig);
        } else if (curLang === 'zh-CN') {
          val = t2s(orig);
        }
        el.setAttribute(attr, val);
      });
    });

    // ③ 對其餘文本節點：
    //    - zh-CN：對所有文本做 t2s（含古典）
    //    - zh-mod（現代漢語）：對所有文本做 t2s；古典區（data-i18n-keep）保留繁體+簡體
    //                          這裡也做 t2s 以保持簡體字面，標題另由 banner 標識為「古文」
    //    - zh-TW：恢復原文
    walkText(root, (textNode) => {
      // 跳過：屬於 data-i18n 元素內部（已在①處理）
      let p = textNode.parentElement;
      while (p) {
        if (p.hasAttribute && p.hasAttribute('data-i18n')) return;
        p = p.parentElement;
      }
      const orig = origText(textNode);
      if (curLang === 'zh-CN' || curLang === 'zh-mod') {
        textNode.nodeValue = t2s(orig);
      } else {
        textNode.nodeValue = orig;
      }
    });

    // ④ 現代漢語模式下：在 data-i18n-keep 區塊頂部插入提示橫幅
    if (curLang === 'zh-mod') {
      const keepNodes = root.querySelectorAll
        ? Array.from(root.querySelectorAll('[data-i18n-keep]'))
        : [];
      keepNodes.forEach(el => {
        if (el.querySelector(':scope > .i18n-classical-banner')) return;
        const banner = document.createElement('div');
        banner.className = 'i18n-classical-banner';
        banner.textContent = '以下为古籍原文，保留古文表达';
        el.insertBefore(banner, el.firstChild);
      });
    } else {
      // 非 EN 模式：移除已存在的橫幅
      const banners = root.querySelectorAll
        ? Array.from(root.querySelectorAll('.i18n-classical-banner'))
        : [];
      banners.forEach(b => b.remove());
    }

    // ⑤ <html lang="..."> 同步
    document.documentElement.setAttribute('lang',
      curLang === 'zh-CN' ? 'zh-CN' : curLang === 'zh-mod' ? 'zh-CN' : 'zh-TW');
  }

  function walkText(root, cb) {
    const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
      acceptNode(n) {
        // 跳過 script/style 內部
        const p = n.parentElement;
        if (!p) return NodeFilter.FILTER_REJECT;
        const tag = p.tagName;
        if (tag === 'SCRIPT' || tag === 'STYLE' || tag === 'NOSCRIPT') return NodeFilter.FILTER_REJECT;
        return NodeFilter.FILTER_ACCEPT;
      }
    });
    let cur;
    while ((cur = walker.nextNode())) cb(cur);
  }

  // ===== 切換 API =====
  function setLang(lang) {
    if (SUPPORTED.indexOf(lang) < 0) return;
    curLang = lang;
    try { localStorage.setItem(STORAGE_KEY, lang); } catch (_) {}
    applyToNode(document.body);
    // 同步切換按鈕高亮
    document.querySelectorAll('[data-lang-switch]').forEach(b => {
      b.classList.toggle('active', b.getAttribute('data-lang-switch') === lang);
    });
    // 派發事件讓其他模塊（location.js / 動態內容渲染）感知
    window.dispatchEvent(new CustomEvent('i18n:changed', { detail: { lang } }));
  }

  // ===== MutationObserver：自動處理動態插入的內容 =====
  function startObserver() {
    if (!('MutationObserver' in window)) return;
    const obs = new MutationObserver(muts => {
      // 每批變更只跑一次 apply，避免抖動
      let touched = false;
      for (const m of muts) {
        if (m.type === 'childList' && m.addedNodes.length) {
          for (const n of m.addedNodes) {
            if (n.nodeType === 1) {
              applyToNode(n);
              touched = true;
            } else if (n.nodeType === 3) {
              // 純文本節點：直接處理
              const orig = origText(n);
              if (curLang === 'zh-CN') n.nodeValue = t2s(orig);
              else n.nodeValue = orig;
              touched = true;
            }
          }
        } else if (m.type === 'characterData') {
          // textContent 被直接賦值的情況：只在 zh-CN 下重新轉換
          // 注意：這裡不能無條件改寫，否則會死循環。我們採用標誌防護。
          if (curLang === 'zh-CN' && !m.target._i18nGuard) {
            m.target._i18nGuard = true;
            const v = m.target.nodeValue;
            const converted = t2s(v);
            if (converted !== v) m.target.nodeValue = converted;
            setTimeout(() => { m.target._i18nGuard = false; }, 0);
          }
        }
      }
      if (touched) {
        // 同步補齊高亮、橫幅等
        document.querySelectorAll('[data-lang-switch]').forEach(b => {
          b.classList.toggle('active', b.getAttribute('data-lang-switch') === curLang);
        });
      }
    });
    obs.observe(document.body, { childList: true, subtree: true, characterData: true });
  }

  // ===== 點擊綁定 =====
  function bindSwitchers() {
    document.querySelectorAll('[data-lang-switch]').forEach(b => {
      b.addEventListener('click', e => {
        e.preventDefault();
        setLang(b.getAttribute('data-lang-switch'));
      });
      b.classList.toggle('active', b.getAttribute('data-lang-switch') === curLang);
    });
  }

  // 啟動：DOM 就緒後做首次轉換 + 綁定 + 監聽
  function boot() {
    applyToNode(document.body);
    bindSwitchers();
    startObserver();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', boot);
  } else {
    boot();
  }

  // 給其他模塊使用：依當前語言返回 key 對應文本
  //   - zh-mod（現代漢語）：查現代白話字典；缺則對 fallback 做繁→簡
  //   - zh-CN：對 fallback 做繁→簡
  //   - zh-TW：直接返回 fallback
  function t(key, fallback) {
    if (curLang === 'zh-mod') {
      const mod = lookupEN(key); // 字典變量名仍叫 EN 但內容已是現代白話
      if (mod) return mod;
      return t2s(fallback || '');
    }
    if (curLang === 'zh-CN') return t2s(fallback || '');
    return fallback || '';
  }

  // 診斷工具：在 console 跑 I18n.debug() 看當前語言、字典條目數、頁面殘留繁體字
  function debug() {
    const dictSize = Object.keys(T2S).length;
    const text = document.body.innerText || document.body.textContent || '';
    const residual = new Map();
    for (const ch of text) {
      const cp = ch.codePointAt(0);
      if (cp >= 0x4E00 && cp <= 0x9FFF && T2S[ch]) {
        residual.set(ch, (residual.get(ch) || 0) + 1);
      }
    }
    const top = Array.from(residual.entries()).sort((a,b) => b[1]-a[1]);
    console.log('[I18n] lang =', curLang, '| dict size =', dictSize);
    console.log('[I18n] residual traditional chars (could be converted but found in body):');
    top.forEach(([c, n]) => console.log(`  '${c}' should be '${T2S[c]}'  ×${n}`));
    if (top.length === 0) console.log('  (none — all converted)');
    return { lang: curLang, dictSize, residual: top };
  }

  window.I18n = { setLang, getLang, t, tr, t2s, apply: applyToNode, debug };
})();
