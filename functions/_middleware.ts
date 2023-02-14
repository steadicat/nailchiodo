const LANGUAGES = ['en', 'it', 'fr'];

function parseCookieHeader(cookieHeader: string | null): { [key: string]: string } {
  if (!cookieHeader) {
    return {};
  }
  const cookies: { [key: string]: string } = {};
  cookieHeader.split(';').forEach(cookie => {
    const [key, value] = cookie.trim().split('=');
    cookies[key] = value;
  });
  return cookies;
}

function getPreferredLanguage(header: string | null, availableLanguages: string[]): string | null {
  if (!header) {
    return null;
  }

  const languages = header
    .split(',')
    .map(language => {
      const [code, q = '1'] = language.trim().split(';q=');
      return { code, q: parseFloat(q) };
    })
    .sort((a, b) => b.q - a.q)
    .map(language => language.code);

  return languages.filter(lang => availableLanguages.includes(lang))[0] ?? null;
}

const languageNegotiation: PagesFunction = async ({ request, next, env }) => {
  const url = new URL(request.url);

  const lang = url.searchParams.get('lang') ?? parseCookieHeader(request.headers.get('cookie'))['lang'] ?? getPreferredLanguage(request.headers.get('accept-language'), LANGUAGES) ?? 'en';

  if (url.pathname.endsWith('/')) {
    url.pathname = `${url.pathname}index.${lang}.html`;
  } else {
    url.pathname = `${url.pathname}.${lang}.html`;
  }

  if (!url.searchParams.get('lang')) {
    return env.ASSETS.fetch(url);
  }

  const asset = await env.ASSETS.fetch(url);
  const response = new Response(asset.body, asset)
  response.headers.append("Set-Cookie", `lang=${lang}; path=/`);
  return response;
};

export const onRequest = [languageNegotiation];