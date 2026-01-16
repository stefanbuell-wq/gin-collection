import { useState, useEffect, useCallback } from 'react';
import { Search, X, Wine, MapPin, Percent, Loader2 } from 'lucide-react';
import { ginReferenceAPI } from '../api/services';
import type { GinReference } from '../types';
import './GinReferenceSearch.css';

interface GinReferenceSearchProps {
  onSelect: (gin: GinReference) => void;
  onClose?: () => void;
  isOpen?: boolean;
}

export function GinReferenceSearch({ onSelect, onClose, isOpen = true }: GinReferenceSearchProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<GinReference[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<{ countries: string[]; gin_types: string[] }>({
    countries: [],
    gin_types: [],
  });
  const [selectedCountry, setSelectedCountry] = useState('');
  const [selectedType, setSelectedType] = useState('');
  const [total, setTotal] = useState(0);

  // Load filters on mount
  useEffect(() => {
    const loadFilters = async () => {
      try {
        const response = await ginReferenceAPI.getFilters();
        if (response.data?.data) {
          setFilters({
            countries: response.data.data.countries || [],
            gin_types: response.data.data.gin_types || [],
          });
        }
      } catch (error) {
        console.error('Failed to load filters:', error);
      }
    };
    loadFilters();
  }, []);

  // Search with debounce
  const searchGins = useCallback(async () => {
    setLoading(true);
    try {
      const response = await ginReferenceAPI.search({
        q: query || undefined,
        country: selectedCountry || undefined,
        type: selectedType || undefined,
        limit: 20,
      });
      if (response.data?.data) {
        setResults(response.data.data.gins || []);
        setTotal(response.data.data.total || 0);
      }
    } catch (error) {
      console.error('Search failed:', error);
      setResults([]);
    } finally {
      setLoading(false);
    }
  }, [query, selectedCountry, selectedType]);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      searchGins();
    }, 300);
    return () => clearTimeout(timer);
  }, [searchGins]);

  // Initial load
  useEffect(() => {
    searchGins();
  }, []);

  const handleSelect = (gin: GinReference) => {
    onSelect(gin);
    if (onClose) onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="gin-reference-search">
      <div className="gin-reference-header">
        <h3>
          <Wine size={20} />
          Gin aus Katalog wählen
        </h3>
        {onClose && (
          <button className="close-btn" onClick={onClose}>
            <X size={20} />
          </button>
        )}
      </div>

      <div className="gin-reference-filters">
        <div className="search-input-wrapper">
          <Search size={18} />
          <input
            type="text"
            placeholder="Suche nach Name, Marke..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            autoFocus
          />
          {query && (
            <button className="clear-btn" onClick={() => setQuery('')}>
              <X size={16} />
            </button>
          )}
        </div>

        <div className="filter-row">
          <select
            value={selectedCountry}
            onChange={(e) => setSelectedCountry(e.target.value)}
          >
            <option value="">Alle Länder</option>
            {filters.countries.map((country) => (
              <option key={country} value={country}>
                {country}
              </option>
            ))}
          </select>

          <select
            value={selectedType}
            onChange={(e) => setSelectedType(e.target.value)}
          >
            <option value="">Alle Typen</option>
            {filters.gin_types.map((type) => (
              <option key={type} value={type}>
                {type}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="gin-reference-results">
        {loading ? (
          <div className="loading-state">
            <Loader2 className="spinner" size={24} />
            <span>Suche...</span>
          </div>
        ) : results.length === 0 ? (
          <div className="empty-state">
            <Wine size={32} />
            <p>Keine Gins gefunden</p>
            <span>Versuche einen anderen Suchbegriff</span>
          </div>
        ) : (
          <>
            <div className="results-count">
              {total} Gins gefunden
            </div>
            <div className="results-list">
              {results.map((gin) => (
                <button
                  key={gin.id}
                  className="gin-reference-item"
                  onClick={() => handleSelect(gin)}
                >
                  <div className="gin-info">
                    <div className="gin-name">{gin.name}</div>
                    <div className="gin-brand">{gin.brand}</div>
                    <div className="gin-meta">
                      {gin.country && (
                        <span className="meta-item">
                          <MapPin size={12} />
                          {gin.country}
                        </span>
                      )}
                      {gin.gin_type && (
                        <span className="meta-item type-badge">
                          {gin.gin_type}
                        </span>
                      )}
                      {gin.abv && (
                        <span className="meta-item">
                          <Percent size={12} />
                          {gin.abv}%
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="select-indicator">
                    Auswählen
                  </div>
                </button>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default GinReferenceSearch;
