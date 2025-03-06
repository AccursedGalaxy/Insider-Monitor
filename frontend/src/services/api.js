/**
 * API client for communication with the backend
 */

const API_BASE_URL = '/api/v1';

/**
 * Generic fetch wrapper with error handling
 */
async function fetchAPI(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || `API request failed with status ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error(`API error for ${url}:`, error);
    throw error;
  }
}

/**
 * Wallet related API calls
 */
export const WalletAPI = {
  // Get all wallets
  getWallets: () => fetchAPI('/wallets'),

  // Get a specific wallet by address
  getWallet: (address) => fetchAPI(`/wallets/${address}`),

  // Get tokens for a specific wallet
  getWalletTokens: (address) => fetchAPI(`/wallets/${address}/tokens`),
};

/**
 * Configuration related API calls
 */
export const ConfigAPI = {
  // Get current configuration
  getConfig: (token) => fetchAPI('/config', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),

  // Update configuration
  updateConfig: (config, token) => fetchAPI('/config', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(config),
  }),

  // Get wallet configurations
  getWalletConfigs: (token) => fetchAPI('/config/wallets', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),

  // Update a wallet configuration
  updateWalletConfig: (address, config, token) => fetchAPI(`/config/wallets/${address}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(config),
  }),

  // Delete a wallet configuration
  deleteWalletConfig: (address, token) => fetchAPI(`/config/wallets/${address}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),
};

/**
 * Alert related API calls
 */
export const AlertAPI = {
  // Get recent alerts
  getAlerts: () => fetchAPI('/alerts'),

  // Get alert settings
  getAlertSettings: (token) => fetchAPI('/alerts/settings', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  }),

  // Update alert settings
  updateAlertSettings: (settings, token) => fetchAPI('/alerts/settings', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(settings),
  }),
};

/**
 * System status related API calls
 */
export const StatusAPI = {
  // Get system status
  getSystemStatus: () => fetchAPI('/status'),

  // Get scan status
  getScanStatus: () => fetchAPI('/status/scan'),
};

/**
 * Health check
 */
export const healthCheck = () => fetch(`${API_BASE_URL}/health`).then(response => response.ok);
