// This would be loaded on the config page

document.addEventListener('DOMContentLoaded', function() {
    // Check if user is logged in
    const token = localStorage.getItem('auth_token');
    if (!token) {
        // Show login form if not logged in
        showLoginForm();
    } else {
        // Load editable configuration if logged in
        loadEditableConfig();
    }

    // Set up login form handler
    document.getElementById('login-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        login();
    });
});

function showLoginForm() {
    const configArea = document.getElementById('config-area');
    if (!configArea) return;

    configArea.innerHTML = `
        <div class="bg-white rounded-lg shadow mb-6">
            <div class="border-b border-gray-200 px-6 py-4">
                <h3 class="text-lg font-medium">Administrator Login</h3>
            </div>
            <div class="p-6">
                <form id="login-form" class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Username</label>
                        <input type="text" id="username" name="username" required
                               class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Password</label>
                        <input type="password" id="password" name="password" required
                               class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                    </div>
                    <div>
                        <button type="submit"
                                class="w-full bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                            Login
                        </button>
                    </div>
                    <div id="login-error" class="text-red-500 text-sm hidden"></div>
                </form>
            </div>
        </div>
    `;

    // Set up form submission
    document.getElementById('login-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        login();
    });
}

function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorEl = document.getElementById('login-error');

    // Clear any previous errors
    if (errorEl) {
        errorEl.classList.add('hidden');
        errorEl.textContent = '';
    }

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Login failed');
        }
        return response.json();
    })
    .then(data => {
        // Store token and load editable config
        localStorage.setItem('auth_token', data.token);
        loadEditableConfig();
    })
    .catch(error => {
        // Show error
        if (errorEl) {
            errorEl.textContent = 'Invalid username or password';
            errorEl.classList.remove('hidden');
        }
    });
}

function loadEditableConfig() {
    const token = localStorage.getItem('auth_token');

    fetch('/api/admin/config', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                // Token expired or invalid
                localStorage.removeItem('auth_token');
                showLoginForm();
                throw new Error('Unauthorized');
            }
            throw new Error('Failed to load config');
        }
        return response.json();
    })
    .then(config => {
        renderEditableConfig(config);
    })
    .catch(error => {
        console.error('Error loading config:', error);
    });
}

function renderEditableConfig(config) {
    const configArea = document.getElementById('config-area');
    if (!configArea) return;

    // Format scan interval for display (remove "ns" suffix if present)
    let scanInterval = config.ScanInterval || "1m";
    if (typeof scanInterval === 'string' && scanInterval.endsWith('ns')) {
        // Convert from nanoseconds string to a human-readable duration
        const ns = parseInt(scanInterval.replace('ns', ''));
        if (ns >= 60000000000) {
            scanInterval = Math.floor(ns / 60000000000) + 'm';
        } else if (ns >= 1000000000) {
            scanInterval = Math.floor(ns / 1000000000) + 's';
        }
    }

    // Ensure scan config exists
    const scanConfig = config.Scan || { ScanMode: 'all', IncludeTokens: [], ExcludeTokens: [] };

    // Ensure wallet configs exists
    const walletConfigs = config.WalletConfigs || {};

    configArea.innerHTML = `
        <div class="bg-white rounded-lg shadow mb-6">
            <div class="border-b border-gray-200 px-6 py-4 flex justify-between">
                <h3 class="text-lg font-medium">Edit Configuration</h3>
                <button id="save-config" class="bg-green-600 text-white px-4 py-1 rounded-md hover:bg-green-700">
                    Save Changes
                </button>
            </div>
            <div class="p-6">
                <form id="config-form" class="space-y-6">
                    <!-- Network Settings -->
                    <div>
                        <h4 class="text-md font-medium mb-2">Network Settings</h4>
                        <div class="bg-gray-50 rounded-lg p-4">
                            <label class="block text-sm font-medium text-gray-700">Network URL</label>
                            <input type="text" id="network-url" value="${config.NetworkURL || ''}"
                                   class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                        </div>
                    </div>

                    <!-- Scan Settings -->
                    <div>
                        <h4 class="text-md font-medium mb-2">Scan Settings</h4>
                        <div class="bg-gray-50 rounded-lg p-4">
                            <div class="mb-4">
                                <label class="block text-sm font-medium text-gray-700">Scan Interval</label>
                                <input type="text" id="scan-interval" value="${scanInterval}"
                                       class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                                <p class="text-xs text-gray-500 mt-1">Format: 1m, 5m, 1h, etc.</p>
                            </div>

                            <div class="border-t border-gray-200 pt-4 mt-4">
                                <label class="block text-sm font-medium text-gray-700 mb-2">Global Scan Mode</label>
                                <select id="scan-mode" class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                                    <option value="all" ${scanConfig.ScanMode === 'all' ? 'selected' : ''}>All Tokens</option>
                                    <option value="whitelist" ${scanConfig.ScanMode === 'whitelist' ? 'selected' : ''}>Whitelist (Only scan included tokens)</option>
                                    <option value="blacklist" ${scanConfig.ScanMode === 'blacklist' ? 'selected' : ''}>Blacklist (Scan all except excluded tokens)</option>
                                </select>
                            </div>

                            <div class="mt-4" id="include-tokens-section" ${scanConfig.ScanMode === 'whitelist' ? '' : 'style="display:none"'}>
                                <label class="block text-sm font-medium text-gray-700 mb-2">Include Tokens (Whitelist)</label>
                                <div class="flex space-x-2 mb-2">
                                    <input type="text" id="new-include-token" placeholder="Token address..."
                                           class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                                    <button type="button" id="add-include-token"
                                            class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                                        Add
                                    </button>
                                </div>
                                <ul id="include-token-list" class="space-y-2 max-h-40 overflow-y-auto">
                                    ${(scanConfig.IncludeTokens || []).map(token => `
                                        <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                            <span class="font-mono">${token}</span>
                                            <button type="button" class="delete-include-token text-red-600 hover:text-red-800" data-token="${token}">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </li>
                                    `).join('')}
                                </ul>
                            </div>

                            <div class="mt-4" id="exclude-tokens-section" ${scanConfig.ScanMode === 'blacklist' ? '' : 'style="display:none"'}>
                                <label class="block text-sm font-medium text-gray-700 mb-2">Exclude Tokens (Blacklist)</label>
                                <div class="flex space-x-2 mb-2">
                                    <input type="text" id="new-exclude-token" placeholder="Token address..."
                                           class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                                    <button type="button" id="add-exclude-token"
                                            class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                                        Add
                                    </button>
                                </div>
                                <ul id="exclude-token-list" class="space-y-2 max-h-40 overflow-y-auto">
                                    ${(scanConfig.ExcludeTokens || []).map(token => `
                                        <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                            <span class="font-mono">${token}</span>
                                            <button type="button" class="delete-exclude-token text-red-600 hover:text-red-800" data-token="${token}">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </li>
                                    `).join('')}
                                </ul>
                            </div>
                        </div>
                    </div>

                    <!-- Wallet Management -->
                    <div>
                        <h4 class="text-md font-medium mb-2">Monitored Wallets</h4>
                        <div class="bg-gray-50 rounded-lg p-4">
                            <div class="mb-4">
                                <div class="flex space-x-2">
                                    <input type="text" id="new-wallet" placeholder="Add a wallet address..."
                                           class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                                    <button type="button" id="add-wallet"
                                            class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                                        Add
                                    </button>
                                </div>
                            </div>

                            <ul id="wallet-list" class="space-y-2 max-h-60 overflow-y-auto">
                                ${(config.Wallets || []).map(wallet => `
                                    <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                        <span class="font-mono">${wallet}</span>
                                        <div class="flex space-x-2">
                                            <button type="button" class="edit-wallet-config text-blue-600 hover:text-blue-800" data-address="${wallet}">
                                                <i class="fas fa-cog" title="Edit wallet-specific settings"></i>
                                            </button>
                                            <button type="button" class="delete-wallet text-red-600 hover:text-red-800" data-address="${wallet}">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </div>
                                    </li>
                                `).join('')}
                            </ul>
                        </div>
                    </div>

                    <!-- Alert Settings -->
                    <div>
                        <h4 class="text-md font-medium mb-2">Alert Settings</h4>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div class="bg-gray-50 rounded-lg p-4">
                                <label class="block text-sm font-medium text-gray-700">Minimum Balance</label>
                                <input type="number" id="min-balance" value="${config.Alerts?.MinimumBalance || 0}"
                                       class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                            </div>

                            <div class="bg-gray-50 rounded-lg p-4">
                                <label class="block text-sm font-medium text-gray-700">Significant Change (%)</label>
                                <input type="number" id="significant-change" value="${(config.Alerts?.SignificantChange || 0) * 100}"
                                       class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                            </div>
                        </div>

                        <div class="mt-4 bg-gray-50 rounded-lg p-4">
                            <label class="block text-sm font-medium text-gray-700 mb-2">Ignored Tokens (for Alerts)</label>
                            <div class="flex space-x-2 mb-2">
                                <input type="text" id="new-ignore-token" placeholder="Token address..."
                                       class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                                <button type="button" id="add-ignore-token"
                                        class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                                    Add
                                </button>
                            </div>
                            <ul id="ignore-token-list" class="space-y-2 max-h-40 overflow-y-auto">
                                ${(config.Alerts?.IgnoreTokens || []).map(token => `
                                    <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                        <span class="font-mono">${token}</span>
                                        <button type="button" class="delete-ignore-token text-red-600 hover:text-red-800" data-token="${token}">
                                            <i class="fas fa-trash"></i>
                                        </button>
                                    </li>
                                `).join('')}
                            </ul>
                        </div>
                    </div>

                    <!-- Discord Settings -->
                    <div>
                        <h4 class="text-md font-medium mb-2">Discord Notifications</h4>
                        <div class="bg-gray-50 rounded-lg p-4">
                            <div class="flex items-center mb-4">
                                <input type="checkbox" id="discord-enabled" ${config.Discord?.Enabled ? 'checked' : ''}
                                       class="h-4 w-4 text-indigo-600 border-gray-300 rounded">
                                <label for="discord-enabled" class="ml-2 block text-sm text-gray-700">
                                    Enable Discord Notifications
                                </label>
                            </div>

                            <div class="space-y-4" id="discord-settings" ${config.Discord?.Enabled ? '' : 'hidden'}>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700">Webhook URL</label>
                                    <input type="text" id="webhook-url" value="${config.Discord?.WebhookURL || ''}"
                                           class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700">Channel ID</label>
                                    <input type="text" id="channel-id" value="${config.Discord?.ChannelID || ''}"
                                           class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                                </div>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>

        <!-- Modal for wallet-specific configuration -->
        <div id="wallet-config-modal" class="hidden fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
            <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl">
                <div class="border-b border-gray-200 px-6 py-4 flex justify-between">
                    <h3 class="text-lg font-medium">Wallet-Specific Settings</h3>
                    <button id="close-wallet-config" class="text-gray-500 hover:text-gray-700">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                <div class="p-6" id="wallet-config-content">
                    <!-- Filled dynamically -->
                </div>
                <div class="border-t border-gray-200 px-6 py-4 flex justify-end space-x-2">
                    <button id="cancel-wallet-config" class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400">
                        Cancel
                    </button>
                    <button id="save-wallet-config" class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700">
                        Save
                    </button>
                </div>
            </div>
        </div>
    `;

    // Set up event handlers
    setupConfigFormHandlers();
}

function setupConfigFormHandlers() {
    // Toggle Discord settings visibility
    const discordEnabled = document.getElementById('discord-enabled');
    const discordSettings = document.getElementById('discord-settings');

    discordEnabled?.addEventListener('change', function() {
        if (discordSettings) {
            discordSettings.hidden = !this.checked;
        }
    });

    // Handle scan mode changes
    const scanMode = document.getElementById('scan-mode');
    scanMode?.addEventListener('change', function() {
        const includeSection = document.getElementById('include-tokens-section');
        const excludeSection = document.getElementById('exclude-tokens-section');

        if (this.value === 'whitelist') {
            includeSection.style.display = '';
            excludeSection.style.display = 'none';
        } else if (this.value === 'blacklist') {
            includeSection.style.display = 'none';
            excludeSection.style.display = '';
        } else { // 'all'
            includeSection.style.display = 'none';
            excludeSection.style.display = 'none';
        }
    });

    // Token management handlers
    document.getElementById('add-include-token')?.addEventListener('click', function() {
        addToken('include');
    });

    document.getElementById('add-exclude-token')?.addEventListener('click', function() {
        addToken('exclude');
    });

    document.getElementById('add-ignore-token')?.addEventListener('click', function() {
        addToken('ignore');
    });

    // Setup delete token handlers
    document.querySelectorAll('.delete-include-token').forEach(btn => {
        btn.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            deleteToken('include', token);
        });
    });

    document.querySelectorAll('.delete-exclude-token').forEach(btn => {
        btn.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            deleteToken('exclude', token);
        });
    });

    document.querySelectorAll('.delete-ignore-token').forEach(btn => {
        btn.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            deleteToken('ignore', token);
        });
    });

    // Add wallet handler
    document.getElementById('add-wallet')?.addEventListener('click', function() {
        addWallet();
    });

    // Delete wallet handlers
    document.querySelectorAll('.delete-wallet').forEach(btn => {
        btn.addEventListener('click', function() {
            const address = this.getAttribute('data-address');
            deleteWallet(address);
        });
    });

    // Edit wallet config handlers
    document.querySelectorAll('.edit-wallet-config').forEach(btn => {
        btn.addEventListener('click', function() {
            const address = this.getAttribute('data-address');
            openWalletConfigModal(address);
        });
    });

    // Modal handlers
    document.getElementById('close-wallet-config')?.addEventListener('click', closeWalletConfigModal);
    document.getElementById('cancel-wallet-config')?.addEventListener('click', closeWalletConfigModal);
    document.getElementById('save-wallet-config')?.addEventListener('click', saveWalletConfig);

    // Save config handler
    document.getElementById('save-config')?.addEventListener('click', function() {
        saveConfig();
    });
}

function addToken(type) {
    const inputId = `new-${type}-token`;
    const listId = `${type}-token-list`;

    const tokenInput = document.getElementById(inputId);
    const tokenList = document.getElementById(listId);

    if (!tokenInput || !tokenList) return;

    const token = tokenInput.value.trim();

    if (!token) {
        alert('Please enter a token address');
        return;
    }

    // Add to UI
    const li = document.createElement('li');
    li.className = 'flex justify-between items-center bg-white p-2 rounded border border-gray-200';
    li.innerHTML = `
        <span class="font-mono">${token}</span>
        <button type="button" class="delete-${type}-token text-red-600 hover:text-red-800" data-token="${token}">
            <i class="fas fa-trash"></i>
        </button>
    `;

    // Add delete handler
    li.querySelector(`.delete-${type}-token`).addEventListener('click', function() {
        deleteToken(type, token);
    });

    tokenList.appendChild(li);

    // Clear input
    tokenInput.value = '';
}

function deleteToken(type, token) {
    // Just remove from UI - will be saved when the config is saved
    const tokenItems = document.querySelectorAll(`#${type}-token-list li`);
    tokenItems.forEach(item => {
        const btn = item.querySelector(`.delete-${type}-token`);
        if (btn && btn.getAttribute('data-token') === token) {
            item.remove();
        }
    });
}

function addWallet() {
    const walletInput = document.getElementById('new-wallet');
    if (!walletInput) return;

    const address = walletInput.value.trim();

    if (!address) {
        alert('Please enter a wallet address');
        return;
    }

    const token = localStorage.getItem('auth_token');

    fetch('/api/admin/wallets', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ address })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to add wallet');
        }
        return response.json();
    })
    .then(data => {
        // Add wallet to the UI
        const walletList = document.getElementById('wallet-list');
        if (walletList) {
            const li = document.createElement('li');
            li.className = 'flex justify-between items-center bg-white p-2 rounded border border-gray-200';
            li.innerHTML = `
                <span class="font-mono">${address}</span>
                <div class="flex space-x-2">
                    <button type="button" class="edit-wallet-config text-blue-600 hover:text-blue-800" data-address="${address}">
                        <i class="fas fa-cog" title="Edit wallet-specific settings"></i>
                    </button>
                    <button type="button" class="delete-wallet text-red-600 hover:text-red-800" data-address="${address}">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            `;

            // Add delete handler to the new button
            li.querySelector('.delete-wallet').addEventListener('click', function() {
                deleteWallet(address);
            });

            walletList.appendChild(li);
        }

        // Clear input
        walletInput.value = '';
    })
    .catch(error => {
        alert('Error adding wallet: ' + error.message);
    });
}

function deleteWallet(address) {
    if (!confirm(`Are you sure you want to remove wallet ${address}?`)) {
        return;
    }

    const token = localStorage.getItem('auth_token');

    fetch(`/api/admin/wallets/${address}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to delete wallet');
        }
        return response.json();
    })
    .then(data => {
        // Remove wallet from the UI
        const walletItems = document.querySelectorAll('#wallet-list li');
        walletItems.forEach(item => {
            const btn = item.querySelector('.delete-wallet');
            if (btn && btn.getAttribute('data-address') === address) {
                item.remove();
            }
        });
    })
    .catch(error => {
        alert('Error deleting wallet: ' + error.message);
    });
}

// New function to open wallet-specific config modal
function openWalletConfigModal(address) {
    const modal = document.getElementById('wallet-config-modal');
    const content = document.getElementById('wallet-config-content');

    if (!modal || !content) return;

    // Get the current configuration
    const token = localStorage.getItem('auth_token');

    fetch('/api/admin/config', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to load config');
        }
        return response.json();
    })
    .then(config => {
        // Get wallet-specific config or create default
        const walletConfig = config.WalletConfigs?.[address] || {};
        const scanConfig = walletConfig.Scan || { ScanMode: 'all', IncludeTokens: [], ExcludeTokens: [] };

        // Set the current wallet address in the modal
        modal.setAttribute('data-wallet', address);

        // Fill the modal content
        content.innerHTML = `
            <p class="text-gray-700 mb-4">Configure wallet-specific settings for: <span class="font-mono font-medium">${address}</span></p>

            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Scan Mode</label>
                    <select id="wallet-scan-mode" class="block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                        <option value="all" ${scanConfig.ScanMode === 'all' ? 'selected' : ''}>All Tokens</option>
                        <option value="whitelist" ${scanConfig.ScanMode === 'whitelist' ? 'selected' : ''}>Whitelist (Only scan included tokens)</option>
                        <option value="blacklist" ${scanConfig.ScanMode === 'blacklist' ? 'selected' : ''}>Blacklist (Scan all except excluded tokens)</option>
                    </select>
                </div>

                <div id="wallet-include-tokens-section" ${scanConfig.ScanMode === 'whitelist' ? '' : 'style="display:none"'}>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Include Tokens (Whitelist)</label>
                    <div class="flex space-x-2 mb-2">
                        <input type="text" id="wallet-new-include-token" placeholder="Token address..."
                               class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                        <button type="button" id="wallet-add-include-token"
                                class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                            Add
                        </button>
                    </div>
                    <ul id="wallet-include-token-list" class="space-y-2 max-h-40 overflow-y-auto">
                        ${(scanConfig.IncludeTokens || []).map(token => `
                            <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                <span class="font-mono">${token}</span>
                                <button type="button" class="wallet-delete-include-token text-red-600 hover:text-red-800" data-token="${token}">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </li>
                        `).join('')}
                    </ul>
                </div>

                <div id="wallet-exclude-tokens-section" ${scanConfig.ScanMode === 'blacklist' ? '' : 'style="display:none"'}>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Exclude Tokens (Blacklist)</label>
                    <div class="flex space-x-2 mb-2">
                        <input type="text" id="wallet-new-exclude-token" placeholder="Token address..."
                               class="flex-1 p-2 border border-gray-300 rounded-md shadow-sm">
                        <button type="button" id="wallet-add-exclude-token"
                                class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">
                            Add
                        </button>
                    </div>
                    <ul id="wallet-exclude-token-list" class="space-y-2 max-h-40 overflow-y-auto">
                        ${(scanConfig.ExcludeTokens || []).map(token => `
                            <li class="flex justify-between items-center bg-white p-2 rounded border border-gray-200">
                                <span class="font-mono">${token}</span>
                                <button type="button" class="wallet-delete-exclude-token text-red-600 hover:text-red-800" data-token="${token}">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </li>
                        `).join('')}
                    </ul>
                </div>
            </div>
        `;

        // Set up wallet modal event handlers
        setupWalletConfigModalHandlers();

        // Show the modal
        modal.classList.remove('hidden');
    })
    .catch(error => {
        alert('Error loading wallet configuration: ' + error.message);
    });
}

function setupWalletConfigModalHandlers() {
    // Handle scan mode changes
    const scanMode = document.getElementById('wallet-scan-mode');
    scanMode?.addEventListener('change', function() {
        const includeSection = document.getElementById('wallet-include-tokens-section');
        const excludeSection = document.getElementById('wallet-exclude-tokens-section');

        if (this.value === 'whitelist') {
            includeSection.style.display = '';
            excludeSection.style.display = 'none';
        } else if (this.value === 'blacklist') {
            includeSection.style.display = 'none';
            excludeSection.style.display = '';
        } else { // 'all'
            includeSection.style.display = 'none';
            excludeSection.style.display = 'none';
        }
    });

    // Token management handlers
    document.getElementById('wallet-add-include-token')?.addEventListener('click', function() {
        addWalletToken('include');
    });

    document.getElementById('wallet-add-exclude-token')?.addEventListener('click', function() {
        addWalletToken('exclude');
    });

    // Setup delete token handlers
    document.querySelectorAll('.wallet-delete-include-token').forEach(btn => {
        btn.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            deleteWalletToken('include', token);
        });
    });

    document.querySelectorAll('.wallet-delete-exclude-token').forEach(btn => {
        btn.addEventListener('click', function() {
            const token = this.getAttribute('data-token');
            deleteWalletToken('exclude', token);
        });
    });
}

function addWalletToken(type) {
    const inputId = `wallet-new-${type}-token`;
    const listId = `wallet-${type}-token-list`;

    const tokenInput = document.getElementById(inputId);
    const tokenList = document.getElementById(listId);

    if (!tokenInput || !tokenList) return;

    const token = tokenInput.value.trim();

    if (!token) {
        alert('Please enter a token address');
        return;
    }

    // Add to UI
    const li = document.createElement('li');
    li.className = 'flex justify-between items-center bg-white p-2 rounded border border-gray-200';
    li.innerHTML = `
        <span class="font-mono">${token}</span>
        <button type="button" class="wallet-delete-${type}-token text-red-600 hover:text-red-800" data-token="${token}">
            <i class="fas fa-trash"></i>
        </button>
    `;

    // Add delete handler
    li.querySelector(`.wallet-delete-${type}-token`).addEventListener('click', function() {
        deleteWalletToken(type, token);
    });

    tokenList.appendChild(li);

    // Clear input
    tokenInput.value = '';
}

function deleteWalletToken(type, token) {
    // Just remove from UI - will be saved when the wallet config is saved
    const tokenItems = document.querySelectorAll(`#wallet-${type}-token-list li`);
    tokenItems.forEach(item => {
        const btn = item.querySelector(`.wallet-delete-${type}-token`);
        if (btn && btn.getAttribute('data-token') === token) {
            item.remove();
        }
    });
}

function closeWalletConfigModal() {
    const modal = document.getElementById('wallet-config-modal');
    if (modal) {
        modal.classList.add('hidden');
    }
}

function saveWalletConfig() {
    const modal = document.getElementById('wallet-config-modal');
    if (!modal) return;

    const walletAddress = modal.getAttribute('data-wallet');

    // Get form values
    const scanMode = document.getElementById('wallet-scan-mode')?.value || 'all';

    // Get token lists
    const includeTokens = [];
    document.querySelectorAll('#wallet-include-token-list li').forEach(item => {
        const tokenText = item.querySelector('span').textContent;
        includeTokens.push(tokenText);
    });

    const excludeTokens = [];
    document.querySelectorAll('#wallet-exclude-token-list li').forEach(item => {
        const tokenText = item.querySelector('span').textContent;
        excludeTokens.push(tokenText);
    });

    // Create the wallet config update object
    const walletConfig = {
        scan: {
            scan_mode: scanMode,
            include_tokens: includeTokens,
            exclude_tokens: excludeTokens
        }
    };

    // Send update to server
    const token = localStorage.getItem('auth_token');

    fetch(`/api/admin/wallet_config/${walletAddress}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(walletConfig)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to update wallet configuration');
        }
        return response.json();
    })
    .then(data => {
        // Show success message and close modal
        alert('Wallet configuration updated successfully');
        closeWalletConfigModal();
    })
    .catch(error => {
        alert('Error updating wallet configuration: ' + error.message);
    });
}

function saveConfig() {
    // Gather form data
    const networkURL = document.getElementById('network-url')?.value;
    const scanInterval = document.getElementById('scan-interval')?.value;
    const minBalance = parseFloat(document.getElementById('min-balance')?.value || '0');
    const significantChange = parseFloat(document.getElementById('significant-change')?.value || '0') / 100;
    const discordEnabled = document.getElementById('discord-enabled')?.checked;
    const webhookURL = document.getElementById('webhook-url')?.value;
    const channelID = document.getElementById('channel-id')?.value;

    // Get scan config
    const scanMode = document.getElementById('scan-mode')?.value || 'all';

    // Get token lists
    const includeTokens = [];
    document.querySelectorAll('#include-token-list li').forEach(item => {
        const tokenText = item.querySelector('span').textContent;
        includeTokens.push(tokenText);
    });

    const excludeTokens = [];
    document.querySelectorAll('#exclude-token-list li').forEach(item => {
        const tokenText = item.querySelector('span').textContent;
        excludeTokens.push(tokenText);
    });

    const ignoreTokens = [];
    document.querySelectorAll('#ignore-token-list li').forEach(item => {
        const tokenText = item.querySelector('span').textContent;
        ignoreTokens.push(tokenText);
    });

    // Create update object with only changed values
    const update = {};

    if (networkURL) update.network_url = networkURL;
    if (scanInterval) update.scan_interval = scanInterval;

    // Add scan config
    update.scan = {
        scan_mode: scanMode,
        include_tokens: includeTokens,
        exclude_tokens: excludeTokens
    };

    // Add alerts
    update.alerts = {
        minimum_balance: minBalance,
        significant_change: significantChange,
        ignore_tokens: ignoreTokens
    };

    // Add discord settings
    update.discord = {
        enabled: discordEnabled
    };

    if (discordEnabled) {
        if (webhookURL) update.discord.webhook_url = webhookURL;
        if (channelID) update.discord.channel_id = channelID;
    }

    // Send update to server
    const token = localStorage.getItem('auth_token');

    fetch('/api/admin/config', {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(update)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to update configuration');
        }
        return response.json();
    })
    .then(data => {
        // Show success message
        alert('Configuration updated successfully');
    })
    .catch(error => {
        alert('Error updating configuration: ' + error.message);
    });
}
