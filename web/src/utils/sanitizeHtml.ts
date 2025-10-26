import DOMPurify from 'dompurify';

const ALLOWED_TAGS = [
  'a',
  'abbr',
  'b',
  'blockquote',
  'br',
  'code',
  'div',
  'em',
  'figure',
  'figcaption',
  'hr',
  'i',
  'img',
  'li',
  'ol',
  'p',
  'span',
  'strong',
  'u',
  'ul',
];

const ALLOWED_ATTR = [
  'href',
  'title',
  'target',
  'rel',
  'src',
  'alt',
  'width',
  'height',
  'loading',
];

const ensureLazyLoading = (html: string) =>
  html.replace(/<img\b(?![^>]*loading=)/gi, '<img loading="lazy" ');

const ensureSafeLinks = (html: string) =>
  html.replace(/<a\b([^>]*)>/gi, (_match, attrs) => {
    const hasTarget = /target=/i.test(attrs);
    const hasRel = /rel=/i.test(attrs);
    const baseAttrs = attrs.trim();
    const parts = [baseAttrs.length > 0 ? baseAttrs : undefined];
    if (!hasTarget) {
      parts.push('target="_blank"');
    }
    if (!hasRel) {
      parts.push('rel="noreferrer"');
    }

    const normalized = parts.filter(Boolean).join(' ').trim();
    return normalized.length > 0 ? `<a ${normalized}>` : '<a>';
  });

export const sanitizeHtml = (value: string) => {
  if (!value) {
    return '';
  }

  const sanitized = DOMPurify.sanitize(value, {
    ALLOWED_TAGS,
    ALLOWED_ATTR,
    ADD_ATTR: ['loading'],
    USE_PROFILES: { html: true },
  });

  return ensureSafeLinks(ensureLazyLoading(sanitized));
};
