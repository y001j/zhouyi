// 访问码客户端：弹窗输入 + localStorage 缓存 + 401 时清缓存
// 输完码即可直接起占，无需再回扉页确认。
(function () {
  'use strict';
  const STORAGE_KEY = 'zhouyi.accessCode';
  const PROTECTED = /^\/api\/(divine|liuren|huican|qimen)\b/;

  function getCode() {
    try { return localStorage.getItem(STORAGE_KEY) || ''; } catch (_) { return ''; }
  }
  function setCode(v) {
    try {
      if (v) localStorage.setItem(STORAGE_KEY, v);
      else localStorage.removeItem(STORAGE_KEY);
    } catch (_) {}
  }
  function clearCode() { setCode(''); }

  // 规整访问码：全角转半角、去空白、转大写、剔除非法字符。
  // 访问码为 8 位 hex（0-9A-F），但放宽到 [A-Z0-9] 容错，避免把合法码误删。
  function sanitizeCode(v) {
    if (!v) return '';
    // 全角 ASCII（U+FF01–U+FF5E）转半角
    let s = String(v).replace(/[！-～]/g, function (ch) {
      return String.fromCharCode(ch.charCodeAt(0) - 0xFEE0);
    });
    // 全角空格转半角，去除所有空白
    s = s.replace(/　/g, ' ').replace(/\s+/g, '');
    s = s.toUpperCase();
    // 只保留 header 安全的 ASCII 字母数字
    s = s.replace(/[^A-Z0-9]/g, '');
    return s;
  }

  function promptCode(message) {
    const def = (window.I18n && window.I18n.t)
      ? window.I18n.t('msg.access.required', '請輸入訪問碼（向管理員索取）')
      : '請輸入訪問碼（向管理員索取）';
    let code = '';
    while (!code) {
      code = window.prompt(message || def, '');
      if (code === null) return ''; // 用户取消
      code = sanitizeCode(code);
      if (!code) continue;
    }
    setCode(code);
    return code;
  }

  async function ensureAccessCode() {
    let code = getCode();
    if (!code) code = promptCode();
    return code;
  }

  // 包装 fetch：受保护接口自动附带 X-Access-Code；401 时清缓存
  async function authedFetch(url, opts) {
    opts = opts || {};
    const headers = new Headers(opts.headers || {});
    const needCode = PROTECTED.test(url);
    if (needCode) {
      let code = sanitizeCode(getCode());
      if (!code) {
        code = promptCode();
        if (!code) throw new Error((window.I18n && window.I18n.t) ? window.I18n.t('msg.access.empty', '未輸入訪問碼') : '未輸入訪問碼');
      }
      headers.set('X-Access-Code', code);
    }
    const res = await fetch(url, Object.assign({}, opts, { headers }));
    if (needCode && res.status === 401) {
      clearCode();
      let body = {};
      try { body = await res.clone().json(); } catch (_) {}
      const fallback = (window.I18n && window.I18n.t) ? window.I18n.t('msg.access.invalid', '訪問碼無效或已使用') : '訪問碼無效或已使用';
      const msg = (body && body.error) ? body.error : fallback;
      const err = new Error(msg);
      err.code = 401;
      throw err;
    }
    // 起卦成功 = 服务端已消费该码：清掉本地，下次需输新码
    if (needCode && res.ok) {
      clearCode();
    }
    return res;
  }

  window.AccessCode = { ensureAccessCode, authedFetch, getCode, clearCode };
})();
