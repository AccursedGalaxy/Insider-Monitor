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
                            <label class="block text-sm font-medium text-gray-700">Scan Interval</label>
                            <input type="text" id="scan-interval" value="${scanInterval}"
                                   class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm">
                            <p class="text-xs text-gray-500 mt-1">Format: 1m, 5m, 1h, etc.</p>
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
                                        <button type="button" class="delete-wallet text-red-600 hover:text-red-800" data-address="${wallet}">
                                            <i class="fas fa-trash"></i>
                                        </button>
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

    // Save config handler
    document.getElementById('save-config')?.addEventListener('click', function() {
        saveConfig();
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
                <button type="button" class="delete-wallet text-red-600 hover:text-red-800" data-address="${address}">
                    <i class="fas fa-trash"></i>
                </button>
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

function saveConfig() {
    // Gather form data
    const networkURL = document.getElementById('network-url')?.value;
    const scanInterval = document.getElementById('scan-interval')?.value;
    const minBalance = parseFloat(document.getElementById('min-balance')?.value || '0');
    const significantChange = parseFloat(document.getElementById('significant-change')?.value || '0') / 100;
    const discordEnabled = document.getElementById('discord-enabled')?.checked;
    const webhookURL = document.getElementById('webhook-url')?.value;
    const channelID = document.getElementById('channel-id')?.value;

    // Create update object with only changed values
    const update = {};

    if (networkURL) update.network_url = networkURL;
    if (scanInterval) update.scan_interval = scanInterval;

    // Add alerts if any alert setting is provided
    if (!isNaN(minBalance) || !isNaN(significantChange)) {
        update.alerts = {};
        if (!isNaN(minBalance)) update.alerts.minimum_balance = minBalance;
        if (!isNaN(significantChange)) update.alerts.significant_change = significantChange;
    }

    // Add discord settings if enabled
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
