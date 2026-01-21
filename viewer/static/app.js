/**
 * Dotfiles Dashboard - Frontend Application
 */

// State
const state = {
    configs: [],
    currentConfig: null,
    originalContent: '',
    editMode: false,
    packages: {
        formulae: [],
        casks: [],
        taps: []
    },
    historyOpen: false,
    historyData: [],
    selectedCommit: null
};

// DOM Elements
const elements = {
    configList: document.getElementById('config-list'),
    welcomeView: document.getElementById('welcome-view'),
    configView: document.getElementById('config-view'),
    packagesView: document.getElementById('packages-view'),
    configTitle: document.getElementById('config-title'),
    configPath: document.getElementById('config-path'),
    configSoftware: document.getElementById('config-software'),
    configFormat: document.getElementById('config-format'),
    configSize: document.getElementById('config-size'),
    configDocs: document.getElementById('config-docs'),
    configRepo: document.getElementById('config-repo'),
    configEditor: document.getElementById('config-editor'),
    codeViewer: document.getElementById('code-viewer'),
    codeContent: document.getElementById('code-content'),
    lineNumbers: document.getElementById('line-numbers'),
    viewModeBtn: document.getElementById('view-mode-btn'),
    editModeBtn: document.getElementById('edit-mode-btn'),
    saveBtn: document.getElementById('save-btn'),
    discardBtn: document.getElementById('discard-btn'),
    downloadBtn: document.getElementById('download-btn'),
    saveStatus: document.getElementById('save-status'),
    reloadBtn: document.getElementById('reload-btn'),
    packagesTitle: document.getElementById('packages-title'),
    packagesList: document.getElementById('packages-list'),
    packageSearch: document.getElementById('package-search'),
    stats: document.getElementById('stats'),
    toast: document.getElementById('toast'),
    historyBtn: document.getElementById('history-btn'),
    historyPanel: document.getElementById('history-panel'),
    historyCloseBtn: document.getElementById('history-close-btn'),
    historyOverlay: document.getElementById('history-overlay'),
    historyList: document.getElementById('history-list'),
    diffView: document.getElementById('diff-view'),
    diffCommitInfo: document.getElementById('diff-commit-info'),
    diffContent: document.getElementById('diff-content')
};

// Syntax Highlighting Rules
const syntaxRules = {
    shell: [
        { pattern: /#.*/g, class: 'hl-comment' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(if|then|else|elif|fi|for|do|done|while|case|esac|function|export|alias|source|local|return|in)\b/g, class: 'hl-keyword' },
        { pattern: /\$\{?[\w]+\}?/g, class: 'hl-variable' },
        { pattern: /\b\d+\b/g, class: 'hl-number' },
    ],
    bash: [
        { pattern: /#.*/g, class: 'hl-comment' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(if|then|else|elif|fi|for|do|done|while|case|esac|function|export|alias|source|local|return|in)\b/g, class: 'hl-keyword' },
        { pattern: /\$\{?[\w]+\}?/g, class: 'hl-variable' },
        { pattern: /\b\d+\b/g, class: 'hl-number' },
    ],
    toml: [
        { pattern: /#.*/g, class: 'hl-comment' },
        { pattern: /^\s*\[[\w\.\-]+\]/gm, class: 'hl-section' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(true|false)\b/g, class: 'hl-boolean' },
        { pattern: /\b\d+\.?\d*\b/g, class: 'hl-number' },
        { pattern: /^[\w\-]+(?=\s*=)/gm, class: 'hl-property' },
    ],
    lua: [
        { pattern: /--\[\[[\s\S]*?\]\]/g, class: 'hl-comment' },
        { pattern: /--.*/g, class: 'hl-comment' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(local|function|end|if|then|else|elseif|for|while|do|return|nil|and|or|not|require|true|false)\b/g, class: 'hl-keyword' },
        { pattern: /\b\d+\.?\d*\b/g, class: 'hl-number' },
    ],
    vim: [
        { pattern: /".*/g, class: 'hl-comment' },
        { pattern: /'[^']*'/g, class: 'hl-string' },
        { pattern: /\b(set|let|if|else|endif|function|endfunction|call|execute|autocmd|augroup|syntax|highlight|map|nmap|imap|vmap|noremap|nnoremap|inoremap|vnoremap)\b/g, class: 'hl-keyword' },
        { pattern: /\b\d+\b/g, class: 'hl-number' },
    ],
    conf: [
        { pattern: /#.*/g, class: 'hl-comment' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(set|bind|unbind|on|off|true|false)\b/gi, class: 'hl-keyword' },
    ],
    ruby: [
        { pattern: /#.*/g, class: 'hl-comment' },
        { pattern: /(["'])(?:\\.|[^\\])*?\1/g, class: 'hl-string' },
        { pattern: /\b(tap|brew|cask|mas|vscode)\b/g, class: 'hl-keyword' },
    ],
};

// Syntax Highlighting
function highlightCode(code, format) {
    const rules = syntaxRules[format] || syntaxRules['conf'];
    let result = escapeHtml(code);

    // Apply highlighting with placeholders to avoid re-matching
    const placeholders = [];

    rules.forEach(rule => {
        result = result.replace(rule.pattern, (match) => {
            const id = `__HL_${placeholders.length}__`;
            placeholders.push(`<span class="${rule.class}">${match}</span>`);
            return id;
        });
    });

    // Replace placeholders with actual spans
    placeholders.forEach((span, i) => {
        result = result.replace(`__HL_${i}__`, span);
    });

    return result;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function generateLineNumbers(content) {
    const lines = content.split('\n');
    return lines.map((_, i) => `<span>${i + 1}</span>`).join('');
}

// API Functions
async function fetchJSON(url, options = {}) {
    const response = await fetch(url, options);
    if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    return response.json();
}

async function loadConfigs() {
    state.configs = await fetchJSON('/api/configs');
    renderConfigList();
}

async function loadConfig(id) {
    const config = await fetchJSON(`/api/configs/${id}`);
    state.currentConfig = config;
    state.originalContent = config.content;
    state.editMode = false;
    renderConfigView(config);
}

async function saveConfig(id, content) {
    return fetchJSON(`/api/configs/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content })
    });
}

async function loadPackages(category) {
    if (!state.packages[category].length) {
        state.packages[category] = await fetchJSON(`/api/brew/${category}`);
    }
    return state.packages[category];
}

async function loadHistory(configId) {
    return fetchJSON(`/api/history/${configId}`);
}

async function loadDiff(configId, commit) {
    return fetchJSON(`/api/history/${configId}/${commit}`);
}

// Render Functions
function renderConfigList() {
    elements.configList.innerHTML = state.configs.map(config => `
        <li>
            <button class="nav-item" data-config-id="${config.id}">
                ${config.displayName}
            </button>
        </li>
    `).join('');

    // Attach click handlers
    elements.configList.querySelectorAll('.nav-item').forEach(item => {
        item.addEventListener('click', () => {
            selectConfig(item.dataset.configId);
        });
    });
}

function renderConfigView(config) {
    // File info
    elements.configTitle.textContent = config.displayName;
    elements.configPath.querySelector('.path-text').textContent = `dotfiles/${config.sourcePath}`;
    elements.configSoftware.textContent = config.software;
    elements.configFormat.textContent = config.format;
    elements.configSize.textContent = formatBytes(config.size);
    elements.configDocs.href = config.docs;
    elements.configRepo.href = config.repo;

    // Set content
    elements.configEditor.value = config.content;

    // Render code viewer with syntax highlighting
    renderCodeViewer(config.content, config.format);

    // Update mode UI
    updateModeUI();
    showView('config');
}

function renderCodeViewer(content, format) {
    // Generate line numbers
    elements.lineNumbers.innerHTML = generateLineNumbers(content);

    // Apply syntax highlighting
    elements.codeContent.innerHTML = highlightCode(content, format);
}

function renderPackagesList(packages, category) {
    const titleMap = {
        formulae: 'Formulae',
        casks: 'Casks',
        taps: 'Taps'
    };

    elements.packagesTitle.textContent = `${titleMap[category]} (${packages.length})`;
    elements.packagesList.innerHTML = packages.map(pkg => `
        <li class="package-item">${pkg.name}</li>
    `).join('');

    showView('packages');
}

function renderStats() {
    const totalConfigs = state.configs.length;

    elements.stats.innerHTML = `
        <div class="stat-card">
            <div class="stat-value">${totalConfigs}</div>
            <div class="stat-label">Config Files</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">${state.packages.formulae.length}</div>
            <div class="stat-label">Formulae</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">${state.packages.casks.length}</div>
            <div class="stat-label">Casks</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">${state.packages.taps.length}</div>
            <div class="stat-label">Taps</div>
        </div>
    `;
}

// History Panel Functions
function toggleHistoryPanel(show) {
    state.historyOpen = show;
    if (show) {
        elements.historyPanel.classList.add('show');
        elements.historyOverlay.classList.add('show');
    } else {
        elements.historyPanel.classList.remove('show');
        elements.historyOverlay.classList.remove('show');
        state.selectedCommit = null;
    }
}

async function openHistoryPanel() {
    if (!state.currentConfig) return;

    toggleHistoryPanel(true);
    elements.historyList.innerHTML = '<div class="history-empty">Loading...</div>';
    elements.diffView.classList.add('hidden');

    try {
        state.historyData = await loadHistory(state.currentConfig.id);
        renderHistoryList();
    } catch (error) {
        elements.historyList.innerHTML = '<div class="history-empty">Failed to load history</div>';
        console.error('History load error:', error);
    }
}

function renderHistoryList() {
    if (state.historyData.length === 0) {
        elements.historyList.innerHTML = '<div class="history-empty">No commit history found</div>';
        return;
    }

    elements.historyList.innerHTML = state.historyData.map(commit => `
        <button class="history-item${state.selectedCommit === commit.fullHash ? ' active' : ''}"
                data-commit="${commit.fullHash}">
            <span class="history-message">
                <span class="history-hash">${commit.hash}</span>
                ${escapeHtml(commit.message)}
            </span>
            <span class="history-meta">${formatDate(commit.date)} by ${escapeHtml(commit.author)}</span>
        </button>
    `).join('');

    // Attach click handlers
    elements.historyList.querySelectorAll('.history-item').forEach(item => {
        item.addEventListener('click', () => {
            selectCommit(item.dataset.commit);
        });
    });
}

async function selectCommit(commitHash) {
    if (!state.currentConfig) return;

    state.selectedCommit = commitHash;

    // Update active state
    elements.historyList.querySelectorAll('.history-item').forEach(item => {
        item.classList.toggle('active', item.dataset.commit === commitHash);
    });

    // Show diff view
    elements.diffView.classList.remove('hidden');
    elements.diffContent.textContent = 'Loading...';

    const commit = state.historyData.find(c => c.fullHash === commitHash);
    if (commit) {
        elements.diffCommitInfo.textContent = `${commit.hash} - ${commit.message}`;
    }

    try {
        const result = await loadDiff(state.currentConfig.id, commitHash);
        renderDiff(result.diff);
    } catch (error) {
        elements.diffContent.textContent = 'Failed to load diff';
        console.error('Diff load error:', error);
    }
}

function renderDiff(diff) {
    if (!diff) {
        elements.diffContent.textContent = 'No changes in this commit';
        return;
    }

    // Apply diff syntax highlighting
    const lines = diff.split('\n');
    const highlighted = lines.map(line => {
        const escaped = escapeHtml(line);
        if (line.startsWith('+') && !line.startsWith('+++')) {
            return `<span class="diff-add">${escaped}</span>`;
        } else if (line.startsWith('-') && !line.startsWith('---')) {
            return `<span class="diff-del">${escaped}</span>`;
        } else if (line.startsWith('@@')) {
            return `<span class="diff-range">${escaped}</span>`;
        } else if (line.startsWith('diff ') || line.startsWith('index ') ||
                   line.startsWith('---') || line.startsWith('+++')) {
            return `<span class="diff-header-line">${escaped}</span>`;
        }
        return escaped;
    }).join('\n');

    elements.diffContent.innerHTML = highlighted;
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    const now = new Date();
    const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;

    return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Mode Management
function updateModeUI() {
    if (state.editMode) {
        // Edit mode
        elements.codeViewer.classList.add('hidden');
        elements.configEditor.classList.remove('hidden');
        elements.viewModeBtn.classList.remove('active');
        elements.editModeBtn.classList.add('active');
        elements.saveBtn.classList.remove('hidden');
        elements.discardBtn.classList.remove('hidden');
        updateEditorButtons();
    } else {
        // View mode
        elements.codeViewer.classList.remove('hidden');
        elements.configEditor.classList.add('hidden');
        elements.viewModeBtn.classList.add('active');
        elements.editModeBtn.classList.remove('active');
        elements.saveBtn.classList.add('hidden');
        elements.discardBtn.classList.add('hidden');
        elements.saveStatus.textContent = '';
    }
}

function toggleEditMode(enable) {
    if (enable && !state.editMode) {
        // Switching to edit mode
        state.editMode = true;
        elements.configEditor.value = state.originalContent;
    } else if (!enable && state.editMode) {
        // Switching to view mode
        if (hasUnsavedChanges()) {
            if (!confirm('Unsaved changes will be lost. Continue?')) {
                return;
            }
        }
        state.editMode = false;
        // Refresh the code viewer
        if (state.currentConfig) {
            renderCodeViewer(state.originalContent, state.currentConfig.format);
        }
    }
    updateModeUI();
}

// View Management
function showView(viewName) {
    elements.welcomeView.classList.remove('active');
    elements.configView.classList.remove('active');
    elements.packagesView.classList.remove('active');

    switch (viewName) {
        case 'welcome':
            elements.welcomeView.classList.add('active');
            break;
        case 'config':
            elements.configView.classList.add('active');
            break;
        case 'packages':
            elements.packagesView.classList.add('active');
            break;
    }
}

function selectConfig(id) {
    // Check for unsaved changes in edit mode
    if (state.editMode && hasUnsavedChanges()) {
        if (!confirm('Unsaved changes will be lost. Continue?')) {
            return;
        }
    }

    // Update active state
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.configId === id) {
            item.classList.add('active');
        }
    });

    loadConfig(id);
}

async function selectPackageCategory(category) {
    // Update active state
    document.querySelectorAll('.nav-expandable').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.category === category) {
            item.classList.add('active');
        }
    });

    document.querySelectorAll('[data-config-id]').forEach(item => {
        item.classList.remove('active');
    });

    const packages = await loadPackages(category);
    renderPackagesList(packages, category);
}

// Editor Functions
function hasUnsavedChanges() {
    return state.editMode &&
           state.currentConfig &&
           elements.configEditor.value !== state.originalContent;
}

function updateEditorButtons() {
    const hasChanges = hasUnsavedChanges();
    elements.saveBtn.disabled = !hasChanges;
    elements.discardBtn.disabled = !hasChanges;

    if (hasChanges) {
        elements.saveStatus.textContent = 'Unsaved changes';
        elements.saveStatus.className = 'status-message';
    } else {
        elements.saveStatus.textContent = '';
    }
}

async function handleSave() {
    if (!state.currentConfig) return;

    const content = elements.configEditor.value;
    elements.saveBtn.disabled = true;
    elements.saveStatus.textContent = 'Saving...';
    elements.saveStatus.className = 'status-message';

    try {
        const result = await saveConfig(state.currentConfig.id, content);
        if (result.success) {
            state.originalContent = content;
            elements.saveStatus.textContent = 'Saved';
            elements.saveStatus.className = 'status-message success';
            showToast('File saved successfully', 'success');

            // Update code viewer with new content
            renderCodeViewer(content, state.currentConfig.format);
        } else {
            throw new Error(result.error || 'Save failed');
        }
    } catch (error) {
        elements.saveStatus.textContent = 'Save failed';
        elements.saveStatus.className = 'status-message error';
        showToast(error.message, 'error');
    }

    updateEditorButtons();
}

function handleDiscard() {
    if (!state.currentConfig) return;

    if (confirm('Discard all changes?')) {
        elements.configEditor.value = state.originalContent;
        updateEditorButtons();
        showToast('Changes discarded', 'success');
    }
}

function handleDownload() {
    if (!state.currentConfig) return;

    const content = state.editMode ? elements.configEditor.value : state.originalContent;
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = state.currentConfig.displayName.replace(/\//g, '_');
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    showToast('Download started', 'success');
}

// Package Search
function handlePackageSearch(query) {
    const items = elements.packagesList.querySelectorAll('.package-item');
    const lowerQuery = query.toLowerCase();

    items.forEach(item => {
        const name = item.textContent.toLowerCase();
        item.style.display = name.includes(lowerQuery) ? '' : 'none';
    });
}

// Utility Functions
function formatBytes(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function showToast(message, type = '') {
    elements.toast.textContent = message;
    elements.toast.className = `toast show ${type}`;

    setTimeout(() => {
        elements.toast.className = 'toast';
    }, 3000);
}

// Event Listeners
function setupEventListeners() {
    // Mode toggle buttons
    elements.viewModeBtn.addEventListener('click', () => toggleEditMode(false));
    elements.editModeBtn.addEventListener('click', () => toggleEditMode(true));

    // Editor events
    elements.configEditor.addEventListener('input', updateEditorButtons);
    elements.saveBtn.addEventListener('click', handleSave);
    elements.discardBtn.addEventListener('click', handleDiscard);
    elements.downloadBtn.addEventListener('click', handleDownload);

    // Reload button
    elements.reloadBtn.addEventListener('click', async () => {
        await init();
        showToast('Reloaded', 'success');
    });

    // Package category buttons
    document.querySelectorAll('.nav-expandable').forEach(btn => {
        btn.addEventListener('click', () => {
            selectPackageCategory(btn.dataset.category);
        });
    });

    // Package search
    elements.packageSearch.addEventListener('input', (e) => {
        handlePackageSearch(e.target.value);
    });

    // History panel
    elements.historyBtn.addEventListener('click', openHistoryPanel);
    elements.historyCloseBtn.addEventListener('click', () => toggleHistoryPanel(false));
    elements.historyOverlay.addEventListener('click', () => toggleHistoryPanel(false));

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        // Ctrl/Cmd + S to save (only in edit mode)
        if ((e.ctrlKey || e.metaKey) && e.key === 's') {
            e.preventDefault();
            if (state.editMode && !elements.saveBtn.disabled) {
                handleSave();
            }
        }
        // Escape to exit edit mode or close history panel
        if (e.key === 'Escape') {
            if (state.historyOpen) {
                toggleHistoryPanel(false);
            } else if (state.editMode) {
                toggleEditMode(false);
            }
        }
    });

    // Warn on page leave with unsaved changes
    window.addEventListener('beforeunload', (e) => {
        if (hasUnsavedChanges()) {
            e.preventDefault();
            e.returnValue = '';
        }
    });
}

// Update package counts in sidebar
async function updatePackageCounts() {
    for (const category of ['formulae', 'casks', 'taps']) {
        const packages = await loadPackages(category);
        const countEl = document.getElementById(`${category}-count`);
        if (countEl) {
            countEl.textContent = packages.length;
        }
    }
}

// Initialize
async function init() {
    try {
        await loadConfigs();
        await updatePackageCounts();
        renderStats();
        showView('welcome');
    } catch (error) {
        console.error('Init error:', error);
        showToast('Failed to load data', 'error');
    }
}

// Setup and start
setupEventListeners();
init();
