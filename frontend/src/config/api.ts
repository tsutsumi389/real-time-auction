/**
 * API設定
 * ローカルネットワーク内での公開に対応するため、動的にURLを構築
 */

/**
 * API Base URLを取得
 * 環境変数が相対パスの場合、またはlocalhost以外からのアクセスの場合に対応
 */
export const getApiBaseUrl = (): string => {
  const envUrl = import.meta.env.VITE_API_BASE_URL;

  // 環境変数が相対パスの場合（例: /api）
  if (envUrl && !envUrl.startsWith('http')) {
    return envUrl;
  }

  // 環境変数が絶対URLの場合
  if (envUrl) {
    return envUrl;
  }

  // デフォルトは相対パス（現在のホストに対するパス）
  return '/api';
};

/**
 * WebSocket URLを取得
 * 現在のプロトコルとホストから動的に構築
 */
export const getWsUrl = (): string => {
  const envUrl = import.meta.env.VITE_WS_URL;

  console.log('[getWsUrl] VITE_WS_URL:', envUrl);
  console.log('[getWsUrl] window.location.protocol:', window.location.protocol);
  console.log('[getWsUrl] window.location.host:', window.location.host);
  console.log('[getWsUrl] window.location.hostname:', window.location.hostname);
  console.log('[getWsUrl] window.location.port:', window.location.port);

  // 環境変数が明示的に設定されている場合（空文字列は無視）
  if (envUrl && envUrl.trim() !== '') {
    console.log('[getWsUrl] Using envUrl:', envUrl);
    return envUrl;
  }

  // 動的にWebSocket URLを構築
  // https の場合は wss、http の場合は ws
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host; // ホスト名とポート番号
  const wsUrl = `${protocol}//${host}/ws`;

  console.log('[getWsUrl] Generated WebSocket URL:', wsUrl);
  return wsUrl;
};

/**
 * 開発用: 設定を出力
 */
export const logApiConfig = (): void => {
  console.log('API Configuration:');
  console.log('  - API Base URL:', getApiBaseUrl());
  console.log('  - WebSocket URL:', getWsUrl());
  console.log('  - Current Host:', window.location.host);
  console.log('  - Current Protocol:', window.location.protocol);
};
