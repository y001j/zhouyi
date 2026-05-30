// 起卦地点（经度）管理模块。
//
// 用于真太阳时校正：钟表时 + (本地经度 − 时区基准经度) × 4 分钟 = 真太阳时。
// 规则：
//   1) localStorage 已缓存 → 直接用
//   2) 否则尝试 navigator.geolocation（HTTPS 才可用），获取后缓存
//   3) 失败/拒绝 → 默认 116.4（北京），并提示用户可在顶部调整
//
// 用法：在业务代码里 await Location.get() 拿到 number；起占请求体里塞 { longitude }
// 同时本模块会自动把一个小标签注入 <body>，点击可手动改经度。
(function () {
  'use strict';
  const STORAGE_KEY = 'zhouyi.longitude';
  const DEFAULT_LON = 116.4; // 北京
  let cached = null;
  let badge = null;

  function readStored() {
    try {
      const v = localStorage.getItem(STORAGE_KEY);
      if (!v) return null;
      const n = parseFloat(v);
      if (isFinite(n) && n >= -180 && n <= 180) return n;
    } catch (_) {}
    return null;
  }

  function store(n) {
    try { localStorage.setItem(STORAGE_KEY, String(n)); } catch (_) {}
  }

  function tryGeo(timeoutMs = 4000) {
    return new Promise(resolve => {
      if (!('geolocation' in navigator)) { resolve(null); return; }
      let done = false;
      const timer = setTimeout(() => { if (!done) { done = true; resolve(null); } }, timeoutMs);
      navigator.geolocation.getCurrentPosition(
        pos => {
          if (done) return;
          done = true; clearTimeout(timer);
          const lon = pos && pos.coords && pos.coords.longitude;
          resolve(typeof lon === 'number' ? lon : null);
        },
        () => { if (!done) { done = true; clearTimeout(timer); resolve(null); } },
        { timeout: timeoutMs, maximumAge: 7 * 24 * 3600 * 1000, enableHighAccuracy: false }
      );
    });
  }

  function fmt(n) { return (Math.round(n * 10) / 10).toFixed(1); }

  function ensureBadge() {
    if (badge) return badge;
    badge = document.createElement('div');
    badge.id = 'locBadge';
    badge.style.cssText = [
      'position:fixed','top:12px','right:12px','z-index:9999',
      'padding:5px 12px','font:12px/1.4 "STKaiti",serif','letter-spacing:.05em',
      'background:rgba(253,246,227,.92)','color:#5a4a2a',
      'border:1px solid #c9a96e','border-radius:14px',
      'cursor:pointer','user-select:none',
      'box-shadow:0 1px 3px rgba(90,60,20,.08)'
    ].join(';');
    badge.title = (window.I18n && window.I18n.t) ? window.I18n.t('loc.title', '點擊修改起卦地經度（用於真太陽時校正）') : '點擊修改起卦地經度（用於真太陽時校正）';
    badge.addEventListener('click', promptManual);
    if (document.body) document.body.appendChild(badge);
    else document.addEventListener('DOMContentLoaded', () => document.body.appendChild(badge));
    return badge;
  }

  function refreshBadge() {
    ensureBadge();
    const t = (window.I18n && window.I18n.t) ? window.I18n.t : (k, f) => f;
    if (cached == null) { badge.textContent = t('loc.notset', '📍 經度未設'); return; }
    const lang = (window.I18n && window.I18n.getLang) ? window.I18n.getLang() : 'zh-TW';
    if (lang === 'en') {
      const ew = cached >= 0 ? 'E' : 'W';
      badge.textContent = `📍 ${ew} ${fmt(Math.abs(cached))}°`;
    } else {
      const ewSrc = cached >= 0 ? '東' : '西';
      const ew = (window.I18n && window.I18n.tr) ? window.I18n.tr(ewSrc) : ewSrc;
      const jingSrc = '經';
      const jing = (window.I18n && window.I18n.tr) ? window.I18n.tr(jingSrc) : jingSrc;
      badge.textContent = `📍 ${ew}${jing} ${fmt(Math.abs(cached))}°`;
    }
  }

  async function promptManual() {
    const cur = cached == null ? DEFAULT_LON : cached;
    const t = (window.I18n && window.I18n.t) ? window.I18n.t : (k, f) => f;
    const v = window.prompt(
      t('loc.prompt', '請輸入起卦地經度（東經為正、西經為負，範圍 -180 ~ 180）：\n例如：北京 116.4，上海 121.5，紐約 -74.0'),
      String(fmt(cur))
    );
    if (v == null) return;
    const n = parseFloat(v);
    if (!isFinite(n) || n < -180 || n > 180) {
      alert(t('loc.invalid', '經度無效，已保留原值'));
      return;
    }
    cached = n;
    store(n);
    refreshBadge();
  }

  // 主入口：返回当前应使用的经度（number）。
  // 总是同步返回一个值，但首次会异步请求 geolocation；获取到会更新 cached 与 badge。
  // 业务代码在每次起卦前 await get() 即可拿到尽量准确的值。
  async function get() {
    if (cached != null) return cached;
    const stored = readStored();
    if (stored != null) {
      cached = stored;
      refreshBadge();
      return cached;
    }
    refreshBadge();
    const geo = await tryGeo();
    if (geo != null) {
      cached = geo;
      store(geo);
    } else {
      cached = DEFAULT_LON; // 默认北京，但不写入 storage——让用户下次访问还能再尝试 geo
    }
    refreshBadge();
    return cached;
  }

  // 页面就绪即触发一次（顺便注入 badge）；业务代码无需等它
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => { get(); });
  } else {
    get();
  }

  // 語言切換時刷新 badge（標籤文字會跟著變）
  window.addEventListener('i18n:changed', () => refreshBadge());

  window.Location = { get, prompt: promptManual };
})();
