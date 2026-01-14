// Main App Logic
class GinCollectionApp {
    constructor() {
        this.apiUrl = 'api/index.php';
        this.currentGins = [];
        this.currentGin = null;
        this.currentView = 'collection';
        
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadGins();
        this.setupRatingInput();
        this.setupPhotoUpload();
    }

    setupEventListeners() {
        // Navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchView(e.target.dataset.view));
        });

        // Search
        document.getElementById('search-btn').addEventListener('click', () => this.searchGins());
        document.getElementById('search-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.searchGins();
        });

        // Filters
        document.getElementById('filter-select').addEventListener('change', () => this.loadGins());
        document.getElementById('sort-select').addEventListener('change', () => this.loadGins());

        // Form
        document.getElementById('gin-form').addEventListener('submit', (e) => this.saveGin(e));
        document.getElementById('cancel-btn').addEventListener('click', () => this.cancelEdit());

        // Scanner
        document.getElementById('scan-btn').addEventListener('click', () => this.openScanner());
        document.getElementById('close-scanner').addEventListener('click', () => this.closeScanner());

        // Detail Modal
        document.getElementById('close-detail').addEventListener('click', () => this.closeDetail());
        document.getElementById('edit-gin-btn').addEventListener('click', () => this.editCurrentGin());
        document.getElementById('delete-gin-btn').addEventListener('click', () => this.deleteCurrentGin());
    }

    switchView(view) {
        // Update navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.view === view);
        });

        // Update views
        document.querySelectorAll('.view').forEach(v => {
            v.classList.toggle('active', v.id === `${view}-view`);
        });

        this.currentView = view;

        // Load data for specific views
        if (view === 'stats') {
            this.loadStats();
        } else if (view === 'add') {
            this.resetForm();
        } else if (view === 'collection') {
            this.loadGins();
        }
    }

    async loadGins() {
        const filter = document.getElementById('filter-select').value;
        const sort = document.getElementById('sort-select').value;

        try {
            const response = await fetch(`${this.apiUrl}?action=list&filter=${filter}&sort=${sort}`);
            const data = await response.json();
            
            this.currentGins = data.gins;
            this.renderGins(data.gins);
        } catch (error) {
            console.error('Error loading gins:', error);
            this.showError('Fehler beim Laden der Gins');
        }
    }

    renderGins(gins) {
        const grid = document.getElementById('gin-grid');
        
        if (gins.length === 0) {
            grid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">🍸</div>
                    <div class="empty-state-text">Noch keine Gins in der Sammlung</div>
                </div>
            `;
            return;
        }

        grid.innerHTML = gins.map(gin => `
            <div class="gin-card" data-id="${gin.id}" onclick="app.showDetail(${gin.id})">
                ${gin.is_finished ? '<div class="gin-card-badge">Leer</div>' : ''}
                <div class="gin-card-image">
                    ${gin.photo_url ? 
                        `<img src="${gin.photo_url}" alt="${gin.name}">` : 
                        '🍸'
                    }
                </div>
                <div class="gin-card-content">
                    <div class="gin-card-title">${gin.name}</div>
                    ${gin.brand ? `<div class="gin-card-brand">${gin.brand}</div>` : ''}
                    <div class="gin-card-info">
                        <div class="gin-card-rating">
                            ${gin.rating ? '⭐'.repeat(gin.rating) : 'Nicht bewertet'}
                        </div>
                        <div class="gin-card-country">${gin.country || ''}</div>
                    </div>
                </div>
            </div>
        `).join('');
    }

    async showDetail(id) {
        try {
            const response = await fetch(`${this.apiUrl}?action=get&id=${id}`);
            const data = await response.json();
            
            if (data.gin) {
                this.currentGin = data.gin;
                this.renderDetail(data.gin);
                document.getElementById('detail-modal').classList.add('active');
            }
        } catch (error) {
            console.error('Error loading gin details:', error);
            this.showError('Fehler beim Laden der Details');
        }
    }

    renderDetail(gin) {
        document.getElementById('detail-name').textContent = gin.name;
        
        const content = document.getElementById('detail-content');
        content.innerHTML = `
            ${gin.photo_url ? `<img src="${gin.photo_url}" alt="${gin.name}" class="detail-photo">` : ''}
            
            <div class="detail-grid">
                ${gin.brand ? `
                    <div class="detail-row">
                        <span class="detail-label">Marke:</span>
                        <span class="detail-value">${gin.brand}</span>
                    </div>
                ` : ''}
                
                ${gin.country ? `
                    <div class="detail-row">
                        <span class="detail-label">Land:</span>
                        <span class="detail-value">${gin.country}</span>
                    </div>
                ` : ''}
                
                ${gin.region ? `
                    <div class="detail-row">
                        <span class="detail-label">Region:</span>
                        <span class="detail-value">${gin.region}</span>
                    </div>
                ` : ''}
                
                ${gin.abv ? `
                    <div class="detail-row">
                        <span class="detail-label">Alkoholgehalt:</span>
                        <span class="detail-value">${gin.abv}%</span>
                    </div>
                ` : ''}
                
                ${gin.bottle_size ? `
                    <div class="detail-row">
                        <span class="detail-label">Flaschengröße:</span>
                        <span class="detail-value">${gin.bottle_size}ml</span>
                    </div>
                ` : ''}
                
                ${gin.price ? `
                    <div class="detail-row">
                        <span class="detail-label">Preis:</span>
                        <span class="detail-value">${gin.price}€</span>
                    </div>
                ` : ''}
                
                ${gin.purchase_date ? `
                    <div class="detail-row">
                        <span class="detail-label">Kaufdatum:</span>
                        <span class="detail-value">${new Date(gin.purchase_date).toLocaleDateString('de-DE')}</span>
                    </div>
                ` : ''}
                
                ${gin.rating ? `
                    <div class="detail-row">
                        <span class="detail-label">Bewertung:</span>
                        <span class="detail-value">${'⭐'.repeat(gin.rating)}</span>
                    </div>
                ` : ''}
                
                ${gin.barcode ? `
                    <div class="detail-row">
                        <span class="detail-label">Barcode:</span>
                        <span class="detail-value">${gin.barcode}</span>
                    </div>
                ` : ''}
            </div>
            
            ${gin.description ? `
                <div class="detail-notes">
                    <strong>Beschreibung:</strong><br>
                    ${gin.description}
                </div>
            ` : ''}
            
            ${gin.tasting_notes ? `
                <div class="detail-notes">
                    <strong>Verkostungsnotizen:</strong><br>
                    ${gin.tasting_notes}
                </div>
            ` : ''}
        `;
    }

    closeDetail() {
        document.getElementById('detail-modal').classList.remove('active');
        this.currentGin = null;
    }

    async searchGins() {
        const query = document.getElementById('search-input').value;
        
        if (!query.trim()) {
            this.loadGins();
            return;
        }

        try {
            const response = await fetch(`${this.apiUrl}?action=search&q=${encodeURIComponent(query)}`);
            const data = await response.json();
            
            this.renderGins(data.gins);
        } catch (error) {
            console.error('Error searching gins:', error);
            this.showError('Fehler bei der Suche');
        }
    }

    async loadStats() {
        try {
            const response = await fetch(`${this.apiUrl}?action=stats`);
            const data = await response.json();
            
            this.renderStats(data.stats);
        } catch (error) {
            console.error('Error loading stats:', error);
            this.showError('Fehler beim Laden der Statistiken');
        }
    }

    renderStats(stats) {
        // Update stat cards
        document.getElementById('stat-total').textContent = stats.total;
        document.getElementById('stat-available').textContent = stats.available;
        document.getElementById('stat-rating').textContent = stats.avg_rating || '-';
        document.getElementById('stat-value').textContent = stats.total_value ? `${stats.total_value}€` : '-';

        // Render countries chart
        const countriesChart = document.getElementById('countries-chart');
        const maxCount = Math.max(...stats.countries.map(c => c.count));
        
        countriesChart.innerHTML = stats.countries.map(country => `
            <div class="chart-bar">
                <div class="chart-label">${country.country}</div>
                <div class="chart-bar-bg">
                    <div class="chart-bar-fill" style="width: ${(country.count / maxCount) * 100}%">
                        ${country.count}
                    </div>
                </div>
            </div>
        `).join('');

        // Render top rated
        const topRatedList = document.getElementById('top-rated-list');
        topRatedList.innerHTML = stats.top_rated.map(gin => `
            <div class="top-rated-item" onclick="app.showDetail(${gin.id})">
                <div class="top-rated-photo">
                    ${gin.photo_url ? `<img src="${gin.photo_url}" alt="${gin.name}">` : '🍸'}
                </div>
                <div class="top-rated-info">
                    <div class="top-rated-name">${gin.name}</div>
                    <div class="top-rated-brand">${gin.brand || ''}</div>
                </div>
                <div class="top-rated-rating">${'⭐'.repeat(gin.rating)}</div>
            </div>
        `).join('');
    }

    setupRatingInput() {
        const stars = document.querySelectorAll('.star');
        const ratingInput = document.getElementById('rating');

        stars.forEach(star => {
            star.addEventListener('click', () => {
                const rating = parseInt(star.dataset.rating);
                ratingInput.value = rating;
                
                stars.forEach(s => {
                    s.classList.toggle('active', parseInt(s.dataset.rating) <= rating);
                    s.textContent = parseInt(s.dataset.rating) <= rating ? '★' : '☆';
                });
            });
        });
    }

    setupPhotoUpload() {
        const photoInput = document.getElementById('photo-input');
        const photoPreview = document.getElementById('photo-preview');

        photoInput.addEventListener('change', async (e) => {
            const file = e.target.files[0];
            if (!file) return;

            // Show preview
            const reader = new FileReader();
            reader.onload = (e) => {
                photoPreview.innerHTML = `<img src="${e.target.result}" alt="Preview">`;
            };
            reader.readAsDataURL(file);

            // Upload file
            const formData = new FormData();
            formData.append('photo', file);

            try {
                const response = await fetch(`${this.apiUrl}?action=upload`, {
                    method: 'POST',
                    body: formData
                });
                const data = await response.json();
                
                if (data.success) {
                    document.getElementById('photo_url').value = data.url;
                }
            } catch (error) {
                console.error('Error uploading photo:', error);
                this.showError('Fehler beim Hochladen des Fotos');
            }
        });
    }

    async saveGin(e) {
        e.preventDefault();

        const formData = {
            name: document.getElementById('name').value,
            brand: document.getElementById('brand').value,
            country: document.getElementById('country').value,
            abv: document.getElementById('abv').value,
            bottle_size: document.getElementById('bottle_size').value,
            price: document.getElementById('price').value,
            purchase_date: document.getElementById('purchase_date').value,
            barcode: document.getElementById('barcode').value,
            rating: document.getElementById('rating').value,
            tasting_notes: document.getElementById('tasting_notes').value,
            description: document.getElementById('description').value,
            photo_url: document.getElementById('photo_url').value
        };

        const ginId = document.getElementById('gin-id').value;
        if (ginId) {
            formData.id = ginId;
        }

        try {
            const action = ginId ? 'update' : 'add';
            const response = await fetch(`${this.apiUrl}?action=${action}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            });
            
            const data = await response.json();
            
            if (data.success) {
                this.resetForm();
                this.switchView('collection');
                this.showSuccess(ginId ? 'Gin aktualisiert!' : 'Gin hinzugefügt!');
            }
        } catch (error) {
            console.error('Error saving gin:', error);
            this.showError('Fehler beim Speichern');
        }
    }

    editCurrentGin() {
        if (!this.currentGin) return;

        // Fill form with current gin data
        document.getElementById('gin-id').value = this.currentGin.id;
        document.getElementById('name').value = this.currentGin.name || '';
        document.getElementById('brand').value = this.currentGin.brand || '';
        document.getElementById('country').value = this.currentGin.country || '';
        document.getElementById('abv').value = this.currentGin.abv || '';
        document.getElementById('bottle_size').value = this.currentGin.bottle_size || 700;
        document.getElementById('price').value = this.currentGin.price || '';
        document.getElementById('purchase_date').value = this.currentGin.purchase_date || '';
        document.getElementById('barcode').value = this.currentGin.barcode || '';
        document.getElementById('rating').value = this.currentGin.rating || '';
        document.getElementById('tasting_notes').value = this.currentGin.tasting_notes || '';
        document.getElementById('description').value = this.currentGin.description || '';
        document.getElementById('photo_url').value = this.currentGin.photo_url || '';

        // Update rating stars
        if (this.currentGin.rating) {
            document.querySelectorAll('.star').forEach(star => {
                const isActive = parseInt(star.dataset.rating) <= this.currentGin.rating;
                star.classList.toggle('active', isActive);
                star.textContent = isActive ? '★' : '☆';
            });
        }

        // Show photo preview if exists
        if (this.currentGin.photo_url) {
            document.getElementById('photo-preview').innerHTML = 
                `<img src="${this.currentGin.photo_url}" alt="Preview">`;
        }

        document.getElementById('form-title').textContent = 'Gin bearbeiten';
        this.closeDetail();
        this.switchView('add');
    }

    async deleteCurrentGin() {
        if (!this.currentGin) return;

        if (!confirm(`Möchtest du "${this.currentGin.name}" wirklich löschen?`)) {
            return;
        }

        try {
            const response = await fetch(`${this.apiUrl}?action=delete`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: this.currentGin.id })
            });
            
            const data = await response.json();
            
            if (data.success) {
                this.closeDetail();
                this.loadGins();
                this.showSuccess('Gin gelöscht!');
            }
        } catch (error) {
            console.error('Error deleting gin:', error);
            this.showError('Fehler beim Löschen');
        }
    }

    resetForm() {
        document.getElementById('gin-form').reset();
        document.getElementById('gin-id').value = '';
        document.getElementById('form-title').textContent = 'Gin hinzufügen';
        document.getElementById('photo-preview').innerHTML = '<span>📸 Foto aufnehmen</span>';
        document.getElementById('purchase_date').value = new Date().toISOString().split('T')[0];
        
        // Reset rating stars
        document.querySelectorAll('.star').forEach(star => {
            star.classList.remove('active');
            star.textContent = '☆';
        });
    }

    cancelEdit() {
        this.resetForm();
        this.switchView('collection');
    }

    openScanner() {
        document.getElementById('scanner-modal').classList.add('active');
        if (window.BarcodeScanner) {
            window.BarcodeScanner.start();
        }
    }

    closeScanner() {
        document.getElementById('scanner-modal').classList.remove('active');
        if (window.BarcodeScanner) {
            window.BarcodeScanner.stop();
        }
    }

    async handleBarcodeScan(barcode) {
        this.closeScanner();
        
        try {
            const response = await fetch(`${this.apiUrl}?action=barcode&code=${barcode}`);
            const data = await response.json();
            
            if (data.exists) {
                // Gin already exists
                if (confirm(`Dieser Gin existiert bereits: "${data.gin.name}". Möchtest du ihn ansehen?`)) {
                    this.showDetail(data.gin.id);
                }
            } else if (data.found) {
                // Found product info, pre-fill form
                document.getElementById('barcode').value = barcode;
                if (data.name) document.getElementById('name').value = data.name;
                if (data.brand) document.getElementById('brand').value = data.brand;
                if (data.country) document.getElementById('country').value = data.country;
                if (data.abv) document.getElementById('abv').value = data.abv;
                
                this.switchView('add');
                this.showSuccess('Produktinfo gefunden!');
            } else {
                // No info found, just set barcode
                document.getElementById('barcode').value = barcode;
                this.switchView('add');
            }
        } catch (error) {
            console.error('Error looking up barcode:', error);
            document.getElementById('barcode').value = barcode;
            this.switchView('add');
        }
    }

    showSuccess(message) {
        // Simple alert for now - can be replaced with toast notifications
        alert(message);
    }

    showError(message) {
        alert('Fehler: ' + message);
    }
}

// Initialize app
const app = new GinCollectionApp();
