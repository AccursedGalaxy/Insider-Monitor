// Common JavaScript functions for Solana Insider Monitor

// Format large numbers with commas
function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

// Format Solana addresses with ellipsis in the middle
function formatAddress(address, prefixLength = 4, suffixLength = 4) {
    if (!address || address.length <= prefixLength + suffixLength) {
        return address;
    }
    return `${address.substring(0, prefixLength)}...${address.substring(address.length - suffixLength)}`;
}

// Convert timestamp to relative time (e.g., "2 hours ago")
function timeAgo(date) {
    if (!(date instanceof Date)) {
        date = new Date(date);
    }

    const seconds = Math.floor((new Date() - date) / 1000);
    let interval = seconds / 31536000;

    if (interval > 1) {
        return Math.floor(interval) + ' years ago';
    }
    interval = seconds / 2592000;
    if (interval > 1) {
        return Math.floor(interval) + ' months ago';
    }
    interval = seconds / 86400;
    if (interval > 1) {
        return Math.floor(interval) + ' days ago';
    }
    interval = seconds / 3600;
    if (interval > 1) {
        return Math.floor(interval) + ' hours ago';
    }
    interval = seconds / 60;
    if (interval > 1) {
        return Math.floor(interval) + ' minutes ago';
    }
    return Math.floor(seconds) + ' seconds ago';
}

// Copy text to clipboard with UI feedback
function copyToClipboard(text, element) {
    navigator.clipboard.writeText(text)
        .then(() => {
            // Provide visual feedback
            const originalHTML = element.innerHTML;
            element.innerHTML = '<i class="fas fa-check text-green-500"></i>';

            setTimeout(() => {
                element.innerHTML = originalHTML;
            }, 2000);
        })
        .catch(err => {
            console.error('Failed to copy text: ', err);
        });
}

// Format percentage with + for positive values
function formatPercentage(value) {
    const sign = value > 0 ? '+' : '';
    return `${sign}${(value * 100).toFixed(2)}%`;
}

// Format token balance based on decimals
function formatTokenBalance(balance, decimals = 0) {
    if (typeof balance !== 'number') {
        return '0';
    }

    const actualBalance = balance / Math.pow(10, decimals);

    // Format based on size
    if (actualBalance >= 1000000) {
        return (actualBalance / 1000000).toFixed(2) + 'M';
    } else if (actualBalance >= 1000) {
        return (actualBalance / 1000).toFixed(2) + 'K';
    } else if (actualBalance < 0.001) {
        return actualBalance.toFixed(8);
    } else {
        return actualBalance.toLocaleString(undefined, {
            minimumFractionDigits: 2,
            maximumFractionDigits: 6
        });
    }
}

// Get network name from URL
function getNetworkName(url) {
    if (!url) return 'Unknown';

    if (url.includes('mainnet')) {
        return 'Mainnet';
    } else if (url.includes('devnet')) {
        return 'Devnet';
    } else if (url.includes('testnet')) {
        return 'Testnet';
    } else {
        return 'Custom RPC';
    }
}

// Calculate percentage change between two values
function calculatePercentageChange(oldValue, newValue) {
    if (oldValue === 0) {
        return newValue > 0 ? 100 : 0;
    }
    return ((newValue - oldValue) / oldValue) * 100;
}

// Theme toggle functionality (light/dark mode)
document.addEventListener('DOMContentLoaded', function() {
    // Check for theme preference in localStorage
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.documentElement.classList.add('dark');
    }

    // Add theme toggle to the page (if button exists)
    const themeToggle = document.getElementById('theme-toggle');
    if (themeToggle) {
        themeToggle.addEventListener('click', function() {
            const isDark = document.documentElement.classList.toggle('dark');
            localStorage.setItem('theme', isDark ? 'dark' : 'light');

            // Update the toggle icon
            const iconElement = themeToggle.querySelector('i');
            if (iconElement) {
                iconElement.className = isDark ? 'fas fa-sun' : 'fas fa-moon';
            }
        });
    }
});
