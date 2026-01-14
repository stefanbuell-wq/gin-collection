// ERWEITERTE FUNKTIONEN für Gin Collection App
// Diese Funktionen zur app.js hinzufügen oder separat einbinden

// Botanicals Management
class BotanicalsManager {
    constructor() {
        this.botanicals = [];
        this.selectedBotanicals = [];
    }

    async loadBotanicals() {
        try {
            const response = await fetch(`${app.apiUrl}?action=botanicals`);
            const data = await response.json();
            this.botanicals = data.botanicals;
            this.renderBotanicalsSelector(data.grouped);
        } catch (error) {
            console.error('Error loading botanicals:', error);
        }
    }

    renderBotanicalsSelector(grouped) {
        const container = document.getElementById('botanicals-selector');
        if (!container) return;

        let html = '';
        for (const [category, items] of Object.entries(grouped)) {
            html += `<div class="botanical-category">${category}</div>`;
            items.forEach(botanical => {
                html += `
                    <div class="botanical-item" data-id="${botanical.id}" 
                         onclick="botanicalsManager.toggleBotanical(${botanical.id}, '${botanical.name}')">
                        ${botanical.name}
                    </div>
                `;
            });
        }
        container.innerHTML = html;
    }

    toggleBotanical(id, name) {
        const element = document.querySelector(`.botanical-item[data-id="${id}"]`);
        const index = this.selectedBotanicals.findIndex(b => b.id === id);

        if (index > -1) {
            this.selectedBotanicals.splice(index, 1);
            element.classList.remove('selected');
        } else {
            this.selectedBotanicals.push({ id, name, prominence: 'notable' });
            element.classList.add('selected');
        }

        document.getElementById('selected_botanicals').value = JSON.stringify(this.selectedBotanicals);
    }

    loadGinBotanicals(ginId) {
        // Bei Bearbeitung: markiere bereits ausgewählte Botanicals
        fetch(`${app.apiUrl}?action=gin-botanicals&gin_id=${ginId}`)
            .then(r => r.json())
            .then(data => {
                this.selectedBotanicals = data.botanicals.map(b => ({
                    id: b.id,
                    name: b.name,
                    prominence: b.prominence
                }));
                
                this.selectedBotanicals.forEach(b => {
                    const element = document.querySelector(`.botanical-item[data-id="${b.id}"]`);
                    if (element) element.classList.add('selected');
                });
            });
    }

    reset() {
        this.selectedBotanicals = [];
        document.querySelectorAll('.botanical-item').forEach(el => {
            el.classList.remove('selected');
        });
    }
}

// Cocktails Manager
class CocktailsManager {
    constructor() {
        this.cocktails = [];
    }

    async loadCocktails() {
        try {
            const response = await fetch(`${app.apiUrl}?action=cocktails`);
            const data = await response.json();
            this.cocktails = data.cocktails;
            return this.cocktails;
        } catch (error) {
            console.error('Error loading cocktails:', error);
            return [];
        }
    }

    async showGinCocktails(ginId) {
        try {
            const response = await fetch(`${app.apiUrl}?action=gin-cocktails&gin_id=${ginId}`);
            const data = await response.json();
            this.renderCocktailsList(data.cocktails);
        } catch (error) {
            console.error('Error loading gin cocktails:', error);
        }
    }

    renderCocktailsList(cocktails) {
        if (cocktails.length === 0) {
            return '<p class="no-cocktails">Keine Cocktail-Empfehlungen verfügbar</p>';
        }

        return cocktails.map(c => `
            <div class="cocktail-card" onclick="cocktailsManager.showCocktail(${c.id})">
                <h4>${c.name}</h4>
                <p>${c.description}</p>
                <div class="cocktail-meta">
                    <span class="difficulty ${c.difficulty}">${this.getDifficultyLabel(c.difficulty)}</span>
                    <span class="prep-time">⏱️ ${c.prep_time} Min</span>
                </div>
            </div>
        `).join('');
    }

    getDifficultyLabel(difficulty) {
        return { easy: 'Einfach', medium: 'Mittel', hard: 'Schwierig' }[difficulty] || difficulty;
    }

    async showCocktail(id) {
        try {
            const response = await fetch(`${app.apiUrl}?action=cocktail&id=${id}`);
            const data = await response.json();
            this.displayCocktailModal(data.cocktail);
        } catch (error) {
            console.error('Error loading cocktail:', error);
        }
    }

    displayCocktailModal(cocktail) {
        const modal = document.createElement('div');
        modal.className = 'modal active';
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>🍸 ${cocktail.name}</h3>
                    <button class="close-btn" onclick="this.closest('.modal').remove()">&times;</button>
                </div>
                <div class="modal-body">
                    <p><strong>Beschreibung:</strong> ${cocktail.description}</p>
                    
                    <h4>Zutaten:</h4>
                    <ul class="ingredients-list">
                        ${cocktail.ingredients.map(i => `
                            <li>${i.amount} ${i.ingredient}</li>
                        `).join('')}
                    </ul>
                    
                    <h4>Zubereitung:</h4>
                    <div class="instructions">${cocktail.instructions.replace(/\n/g, '<br>')}</div>
                    
                    <div class="cocktail-details">
                        <span><strong>Glas:</strong> ${cocktail.glass_type}</span>
                        <span><strong>Eis:</strong> ${cocktail.ice_type}</span>
                        <span><strong>Schwierigkeit:</strong> ${this.getDifficultyLabel(cocktail.difficulty)}</span>
                        <span><strong>Zeit:</strong> ${cocktail.prep_time} Min</span>
                    </div>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
    }
}

// AI Suggestions
class AISuggestionsManager {
    async loadSuggestions(ginId) {
        try {
            const response = await fetch(`${app.apiUrl}?action=ai-suggestions&gin_id=${ginId}`);
            const data = await response.json();
            return data.suggestions;
        } catch (error) {
            console.error('Error loading AI suggestions:', error);
            return {};
        }
    }

    renderSuggestions(suggestions) {
        if (!suggestions) return '';

        let html = '<div class="ai-suggestions">';
        html += '<h3>💡 Du könntest auch mögen</h3>';

        if (suggestions.same_country && suggestions.same_country.length > 0) {
            html += '<div class="suggestion-section">';
            html += '<h4>Aus dem gleichen Land</h4>';
            html += '<div class="suggestion-grid">';
            suggestions.same_country.forEach(gin => {
                html += this.renderSuggestionCard(gin);
            });
            html += '</div></div>';
        }

        if (suggestions.shared_botanicals && suggestions.shared_botanicals.length > 0) {
            html += '<div class="suggestion-section">';
            html += '<h4>Ähnliche Botanicals</h4>';
            html += '<div class="suggestion-grid">';
            suggestions.shared_botanicals.forEach(gin => {
                html += this.renderSuggestionCard(gin);
            });
            html += '</div></div>';
        }

        html += '</div>';
        return html;
    }

    renderSuggestionCard(gin) {
        return `
            <div class="suggestion-card" onclick="app.showDetail(${gin.id})">
                <div class="suggestion-photo">
                    ${gin.photo_url ? 
                        `<img src="${gin.photo_url}" alt="${gin.name}">` : 
                        '🍸'
                    }
                </div>
                <div class="suggestion-info">
                    <div class="suggestion-name">${gin.name}</div>
                    <div class="suggestion-brand">${gin.brand || ''}</div>
                    ${gin.rating ? `<div class="suggestion-rating">${'⭐'.repeat(gin.rating)}</div>` : ''}
                </div>
            </div>
        `;
    }
}

// Export/Import Manager
class DataManager {
    async exportJSON() {
        try {
            const response = await fetch(`${app.apiUrl}?action=export&format=json`);
            const data = await response.json();
            
            const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `gin-collection-${data.export_date}.json`;
            a.click();
            
            app.showSuccess(`${data.total_gins} Gins exportiert!`);
        } catch (error) {
            console.error('Export error:', error);
            app.showError('Export fehlgeschlagen');
        }
    }

    async exportCSV() {
        window.open(`${app.apiUrl}?action=export&format=csv`, '_blank');
    }

    async importJSON(file) {
        try {
            const text = await file.text();
            const data = JSON.parse(text);
            
            const response = await fetch(`${app.apiUrl}?action=import`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            
            const result = await response.json();
            
            if (result.success) {
                app.showSuccess(`${result.imported} Gins importiert!`);
                if (result.errors.length > 0) {
                    console.warn('Import errors:', result.errors);
                }
                app.loadGins();
            }
        } catch (error) {
            console.error('Import error:', error);
            app.showError('Import fehlgeschlagen');
        }
    }
}

// Füllstand-Update
function updateFillLevel() {
    const slider = document.getElementById('fill_level');
    const display = document.getElementById('fill_level_display');
    const bar = document.getElementById('fill_level_bar');
    
    if (slider && display && bar) {
        slider.addEventListener('input', (e) => {
            const value = e.target.value;
            display.textContent = value + '%';
            bar.style.width = value + '%';
            
            // Farbe je nach Füllstand
            if (value <= 25) {
                bar.style.background = 'linear-gradient(90deg, #e74c3c, #c0392b)';
            } else if (value <= 50) {
                bar.style.background = 'linear-gradient(90deg, #f39c12, #e67e22)';
            } else if (value <= 75) {
                bar.style.background = 'linear-gradient(90deg, #f1c40f, #f39c12)';
            } else {
                bar.style.background = 'linear-gradient(90deg, #27ae60, #2ecc71)';
            }
        });
    }
}

// Globale Instanzen
const botanicalsManager = new BotanicalsManager();
const cocktailsManager = new CocktailsManager();
const aiSuggestionsManager = new AISuggestionsManager();
const dataManager = new DataManager();

// Initialisierung
document.addEventListener('DOMContentLoaded', () => {
    botanicalsManager.loadBotanicals();
    updateFillLevel();
});

// CSS für neue Features
const extendedStyles = `
<style>
.cocktail-card {
    background: white;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
    cursor: pointer;
    transition: transform 0.2s;
}

.cocktail-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-hover);
}

.cocktail-meta {
    display: flex;
    gap: 1rem;
    margin-top: 0.5rem;
    font-size: 0.9rem;
}

.difficulty {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
}

.difficulty.easy { background: #27ae60; color: white; }
.difficulty.medium { background: #f39c12; color: white; }
.difficulty.hard { background: #e74c3c; color: white; }

.ingredients-list {
    list-style: none;
    padding: 0;
}

.ingredients-list li {
    padding: 0.5rem;
    background: var(--bg-color);
    margin-bottom: 0.5rem;
    border-radius: 4px;
}

.instructions {
    background: var(--bg-color);
    padding: 1rem;
    border-radius: 8px;
    line-height: 1.6;
}

.cocktail-details {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border-color);
}

.ai-suggestions {
    margin-top: 2rem;
    padding: 1.5rem;
    background: linear-gradient(135deg, #667eea22, #764ba222);
    border-radius: 10px;
}

.suggestion-section {
    margin-top: 1.5rem;
}

.suggestion-section h4 {
    color: var(--primary-color);
    margin-bottom: 1rem;
}

.suggestion-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 1rem;
}

.suggestion-card {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.2s;
}

.suggestion-card:hover {
    transform: translateY(-3px);
    box-shadow: var(--shadow-hover);
}

.suggestion-photo {
    width: 100%;
    height: 120px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2.5rem;
}

.suggestion-photo img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.suggestion-info {
    padding: 0.75rem;
}

.suggestion-name {
    font-weight: 600;
    font-size: 0.9rem;
    color: var(--primary-color);
}

.suggestion-brand {
    font-size: 0.8rem;
    color: var(--text-light);
}

.suggestion-rating {
    font-size: 0.8rem;
    margin-top: 0.25rem;
}

.no-cocktails {
    text-align: center;
    color: var(--text-light);
    padding: 2rem;
}
</style>
`;
