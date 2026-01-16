import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useGinStore } from '../stores/ginStore';
import {
  Wine,
  ArrowLeft,
  Save,
  Star,
  MapPin,
  Calendar,
  Euro,
  Percent,
  Droplets,
  FileText,
  Loader2,
  AlertCircle,
  Barcode,
  Globe,
  Tag,
  Beaker,
  ShoppingBag,
  Sparkles,
  Cherry,
  BookOpen
} from 'lucide-react';
import type { GinCreateRequest, GinReference } from '../types';
import { GinReferenceSearch } from '../components/GinReferenceSearch';
import './GinCreate.css';

// Common gin types for dropdown
const GIN_TYPES = [
  'London Dry',
  'Plymouth',
  'Old Tom',
  'Genever',
  'New Western',
  'Navy Strength',
  'Sloe Gin',
  'Flavoured Gin',
  'Other'
];

// Common countries for dropdown
const COUNTRIES = [
  'Deutschland',
  'Vereinigtes Königreich',
  'Schottland',
  'Spanien',
  'Niederlande',
  'Belgien',
  'Frankreich',
  'Italien',
  'USA',
  'Japan',
  'Australien',
  'Andere'
];

// Common bottle sizes
const BOTTLE_SIZES = [50, 100, 200, 350, 500, 700, 750, 1000, 1500];

const GinCreate = () => {
  const navigate = useNavigate();
  const { createGin, isLoading, error, clearError } = useGinStore();

  const [formData, setFormData] = useState<GinCreateRequest>({
    name: '',
    brand: '',
    country: '',
    region: '',
    gin_type: '',
    abv: undefined,
    bottle_size: 700,
    fill_level: 100,
    price: undefined,
    current_market_value: undefined,
    barcode: '',
    purchase_date: '',
    purchase_location: '',
    rating: undefined,
    nose_notes: '',
    palate_notes: '',
    finish_notes: '',
    general_notes: '',
    description: '',
    recommended_tonic: '',
    recommended_garnish: '',
    is_finished: false
  });

  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  const [showCatalog, setShowCatalog] = useState(false);

  // Handle selection from gin catalog
  const handleCatalogSelect = (gin: GinReference) => {
    setFormData(prev => ({
      ...prev,
      name: gin.name,
      brand: gin.brand || '',
      country: gin.country || '',
      region: gin.region || '',
      gin_type: gin.gin_type || '',
      abv: gin.abv,
      bottle_size: gin.bottle_size || 700,
      description: gin.description || '',
      nose_notes: gin.nose_notes || '',
      palate_notes: gin.palate_notes || '',
      finish_notes: gin.finish_notes || '',
      recommended_tonic: gin.recommended_tonic || '',
      recommended_garnish: gin.recommended_garnish || '',
      barcode: gin.barcode || '',
    }));
    setShowCatalog(false);
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;

    // Clear validation error when field changes
    if (validationErrors[name]) {
      setValidationErrors(prev => {
        const next = { ...prev };
        delete next[name];
        return next;
      });
    }

    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData(prev => ({ ...prev, [name]: checked }));
    } else if (type === 'number') {
      setFormData(prev => ({
        ...prev,
        [name]: value === '' ? undefined : parseFloat(value)
      }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  const handleRatingClick = (rating: number) => {
    setFormData(prev => ({
      ...prev,
      rating: prev.rating === rating ? undefined : rating
    }));
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.name?.trim()) {
      errors.name = 'Name ist erforderlich';
    }

    if (formData.abv !== undefined && (formData.abv < 0 || formData.abv > 100)) {
      errors.abv = 'ABV muss zwischen 0 und 100 liegen';
    }

    if (formData.fill_level !== undefined && (formData.fill_level < 0 || formData.fill_level > 100)) {
      errors.fill_level = 'Füllstand muss zwischen 0 und 100 liegen';
    }

    if (formData.price !== undefined && formData.price < 0) {
      errors.price = 'Preis kann nicht negativ sein';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();

    if (!validate()) {
      return;
    }

    try {
      const ginData = { ...formData };
      // Clean up empty strings
      Object.keys(ginData).forEach(key => {
        const k = key as keyof GinCreateRequest;
        if (ginData[k] === '') {
          (ginData as any)[k] = undefined;
        }
      });

      await createGin(ginData);
      navigate('/gins');
    } catch (err) {
      console.error('Failed to create gin:', err);
    }
  };

  // Animation variants
  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 15
      }
    }
  };

  // Render stars
  const renderStars = () => {
    return [1, 2, 3, 4, 5].map((star) => (
      <motion.button
        key={star}
        type="button"
        onClick={() => handleRatingClick(star)}
        className={`gin-create-star ${formData.rating && star <= formData.rating ? 'gin-create-star--filled' : ''}`}
        whileHover={{ scale: 1.15 }}
        whileTap={{ scale: 0.95 }}
      >
        <Star />
      </motion.button>
    ));
  };

  return (
    <div className="gin-create-page">
      {/* Ambient Background */}
      <div className="gin-create-ambient">
        <div className="gin-create-orb gin-create-orb--gold" />
        <div className="gin-create-orb gin-create-orb--mint" />
        <div className="gin-create-orb gin-create-orb--green" />
      </div>

      <div className="gin-create-content">
        {/* Header */}
        <motion.div
          className="gin-create-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <motion.button
            onClick={() => navigate('/gins')}
            className="gin-create-back-btn"
            whileHover={{ x: -4 }}
            whileTap={{ scale: 0.95 }}
          >
            <ArrowLeft />
          </motion.button>

          <div className="gin-create-header__info">
            <h1 className="gin-create-header__title">
              <div className="gin-create-header__icon">
                <Wine />
              </div>
              Neuen Gin hinzufügen
            </h1>
            <p className="gin-create-header__subtitle">
              Füge eine neue Flasche zu deinem Tresor hinzu
            </p>
          </div>

          <motion.button
            type="button"
            onClick={() => setShowCatalog(true)}
            className="gin-create-catalog-btn"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <BookOpen />
            Aus Katalog wählen
          </motion.button>
        </motion.div>

        {/* Gin Catalog Modal */}
        <AnimatePresence>
          {showCatalog && (
            <motion.div
              className="gin-reference-modal-overlay"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={(e) => {
                if (e.target === e.currentTarget) setShowCatalog(false);
              }}
            >
              <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 20 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 20 }}
                transition={{ type: 'spring', damping: 25 }}
              >
                <GinReferenceSearch
                  onSelect={handleCatalogSelect}
                  onClose={() => setShowCatalog(false)}
                  isOpen={true}
                />
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Error Display */}
        {error && (
          <motion.div
            className="gin-create-error"
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <AlertCircle className="gin-create-error__icon" />
            <span className="gin-create-error__text">{error}</span>
          </motion.div>
        )}

        <form onSubmit={handleSubmit}>
          {/* Basic Information */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.1 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon">
                <Wine />
              </div>
              <h2 className="gin-create-card__title">Grundinformationen</h2>
            </div>

            <div className="gin-create-grid">
              <div className="gin-create-group gin-create-grid__full">
                <label className="gin-create-label">
                  Name <span className="gin-create-label__required">*</span>
                </label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="z.B. Monkey 47 Schwarzwald Dry Gin"
                  className={`gin-create-input ${validationErrors.name ? 'gin-create-input--error' : ''}`}
                  required
                />
                {validationErrors.name && (
                  <span className="gin-create-field-error">
                    <AlertCircle />
                    {validationErrors.name}
                  </span>
                )}
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Marke</label>
                <input
                  type="text"
                  name="brand"
                  value={formData.brand}
                  onChange={handleChange}
                  placeholder="z.B. Monkey 47"
                  className="gin-create-input"
                />
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Gin-Typ</label>
                <select
                  name="gin_type"
                  value={formData.gin_type}
                  onChange={handleChange}
                  className="gin-create-select"
                >
                  <option value="">Typ auswählen...</option>
                  {GIN_TYPES.map(type => (
                    <option key={type} value={type}>{type}</option>
                  ))}
                </select>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Land</label>
                <select
                  name="country"
                  value={formData.country}
                  onChange={handleChange}
                  className="gin-create-select"
                >
                  <option value="">Land auswählen...</option>
                  {COUNTRIES.map(country => (
                    <option key={country} value={country}>{country}</option>
                  ))}
                </select>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Region</label>
                <input
                  type="text"
                  name="region"
                  value={formData.region}
                  onChange={handleChange}
                  placeholder="z.B. Schwarzwald"
                  className="gin-create-input"
                />
              </div>

              <div className="gin-create-group gin-create-grid__full">
                <label className="gin-create-label">Beschreibung</label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  placeholder="Kurze Beschreibung dieses Gins..."
                  rows={2}
                  className="gin-create-textarea"
                />
              </div>
            </div>
          </motion.div>

          {/* Technical Details */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.2 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon gin-create-card__icon--mint">
                <Beaker />
              </div>
              <h2 className="gin-create-card__title">Technische Details</h2>
            </div>

            <div className="gin-create-grid gin-create-grid--thirds">
              <div className="gin-create-group">
                <label className="gin-create-label">ABV (%)</label>
                <input
                  type="number"
                  name="abv"
                  value={formData.abv ?? ''}
                  onChange={handleChange}
                  placeholder="z.B. 47"
                  min="0"
                  max="100"
                  step="0.1"
                  className={`gin-create-input ${validationErrors.abv ? 'gin-create-input--error' : ''}`}
                />
                {validationErrors.abv && (
                  <span className="gin-create-field-error">
                    <AlertCircle />
                    {validationErrors.abv}
                  </span>
                )}
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Flaschengröße (ml)</label>
                <select
                  name="bottle_size"
                  value={formData.bottle_size ?? ''}
                  onChange={handleChange}
                  className="gin-create-select"
                >
                  <option value="">Größe wählen...</option>
                  {BOTTLE_SIZES.map(size => (
                    <option key={size} value={size}>{size} ml</option>
                  ))}
                </select>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Füllstand</label>
                <div className="gin-create-slider-container">
                  <input
                    type="range"
                    name="fill_level"
                    value={formData.fill_level ?? 100}
                    onChange={handleChange}
                    min="0"
                    max="100"
                    className="gin-create-slider"
                  />
                  <span className="gin-create-slider-value">
                    {formData.fill_level ?? 100}%
                  </span>
                </div>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Barcode / EAN</label>
                <div className="gin-create-input-wrapper">
                  <input
                    type="text"
                    name="barcode"
                    value={formData.barcode}
                    onChange={handleChange}
                    placeholder="z.B. 4260145970016"
                    className="gin-create-input gin-create-input--barcode"
                  />
                  <motion.button
                    type="button"
                    className="gin-create-barcode-btn"
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    title="Barcode scannen"
                  >
                    <Barcode />
                  </motion.button>
                </div>
              </div>

              <div className="gin-create-group gin-create-grid__full">
                <label className="gin-create-checkbox-label">
                  <input
                    type="checkbox"
                    name="is_finished"
                    checked={formData.is_finished}
                    onChange={handleChange}
                    className="gin-create-checkbox"
                  />
                  <span className="gin-create-checkbox-text">
                    Flasche ist leer (aufgebraucht)
                  </span>
                </label>
              </div>
            </div>
          </motion.div>

          {/* Purchase Information */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.3 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon">
                <ShoppingBag />
              </div>
              <h2 className="gin-create-card__title">Kaufinformationen</h2>
            </div>

            <div className="gin-create-grid">
              <div className="gin-create-group">
                <label className="gin-create-label">Kaufpreis</label>
                <div className="gin-create-input-wrapper">
                  <span className="gin-create-input-prefix">EUR</span>
                  <input
                    type="number"
                    name="price"
                    value={formData.price ?? ''}
                    onChange={handleChange}
                    placeholder="0.00"
                    min="0"
                    step="0.01"
                    className={`gin-create-input gin-create-input--prefix ${validationErrors.price ? 'gin-create-input--error' : ''}`}
                  />
                </div>
                {validationErrors.price && (
                  <span className="gin-create-field-error">
                    <AlertCircle />
                    {validationErrors.price}
                  </span>
                )}
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Aktueller Marktwert</label>
                <div className="gin-create-input-wrapper">
                  <span className="gin-create-input-prefix">EUR</span>
                  <input
                    type="number"
                    name="current_market_value"
                    value={formData.current_market_value ?? ''}
                    onChange={handleChange}
                    placeholder="0.00"
                    min="0"
                    step="0.01"
                    className="gin-create-input gin-create-input--prefix"
                  />
                </div>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Kaufdatum</label>
                <div className="gin-create-input-wrapper gin-create-input-wrapper--icon">
                  <input
                    type="date"
                    name="purchase_date"
                    value={formData.purchase_date}
                    onChange={handleChange}
                    className="gin-create-input"
                  />
                  <div className="gin-create-input-icon">
                    <Calendar />
                  </div>
                </div>
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Kaufort</label>
                <div className="gin-create-input-wrapper gin-create-input-wrapper--icon">
                  <input
                    type="text"
                    name="purchase_location"
                    value={formData.purchase_location}
                    onChange={handleChange}
                    placeholder="z.B. Amazon, Lokal"
                    className="gin-create-input"
                  />
                  <div className="gin-create-input-icon">
                    <MapPin />
                  </div>
                </div>
              </div>
            </div>
          </motion.div>

          {/* Rating */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.4 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon">
                <Sparkles />
              </div>
              <h2 className="gin-create-card__title">Bewertung</h2>
            </div>

            <div className="gin-create-rating">
              {renderStars()}
              {formData.rating && (
                <span className="gin-create-rating-value">
                  {formData.rating}/5
                </span>
              )}
            </div>
          </motion.div>

          {/* Tasting Notes */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.5 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon">
                <FileText />
              </div>
              <h2 className="gin-create-card__title">Verkostungsnotizen</h2>
            </div>

            <div className="gin-create-tasting-grid">
              <div className="gin-create-tasting-item">
                <label className="gin-create-tasting-label">Nase</label>
                <textarea
                  name="nose_notes"
                  value={formData.nose_notes}
                  onChange={handleChange}
                  placeholder="Welche Aromen riechst du? z.B. Wacholder, Zitrus, Blumen..."
                  rows={2}
                  className="gin-create-textarea"
                />
              </div>

              <div className="gin-create-tasting-item">
                <label className="gin-create-tasting-label">Gaumen</label>
                <textarea
                  name="palate_notes"
                  value={formData.palate_notes}
                  onChange={handleChange}
                  placeholder="Welche Geschmäcker schmeckst du? z.B. Wacholder, Pfeffer, Süße..."
                  rows={2}
                  className="gin-create-textarea"
                />
              </div>

              <div className="gin-create-tasting-item">
                <label className="gin-create-tasting-label">Abgang</label>
                <textarea
                  name="finish_notes"
                  value={formData.finish_notes}
                  onChange={handleChange}
                  placeholder="Wie ist der Abgang? z.B. Lang, warm, Zitrus nachklingend..."
                  rows={2}
                  className="gin-create-textarea"
                />
              </div>

              <div className="gin-create-tasting-item">
                <label className="gin-create-tasting-label">Allgemeine Notizen</label>
                <textarea
                  name="general_notes"
                  value={formData.general_notes}
                  onChange={handleChange}
                  placeholder="Weitere Notizen zu diesem Gin..."
                  rows={2}
                  className="gin-create-textarea"
                />
              </div>
            </div>
          </motion.div>

          {/* Serving Suggestions */}
          <motion.div
            className="gin-create-card"
            variants={cardVariants}
            initial="hidden"
            animate="visible"
            transition={{ delay: 0.6 }}
          >
            <div className="gin-create-card__header">
              <div className="gin-create-card__icon gin-create-card__icon--mint">
                <Cherry />
              </div>
              <h2 className="gin-create-card__title">Serviervorschläge</h2>
            </div>

            <div className="gin-create-grid">
              <div className="gin-create-group">
                <label className="gin-create-label">Empfohlenes Tonic</label>
                <input
                  type="text"
                  name="recommended_tonic"
                  value={formData.recommended_tonic}
                  onChange={handleChange}
                  placeholder="z.B. Fever-Tree Mediterranean"
                  className="gin-create-input"
                />
              </div>

              <div className="gin-create-group">
                <label className="gin-create-label">Empfohlene Garnitur</label>
                <input
                  type="text"
                  name="recommended_garnish"
                  value={formData.recommended_garnish}
                  onChange={handleChange}
                  placeholder="z.B. Grapefruit-Schale, Rosmarin"
                  className="gin-create-input"
                />
              </div>
            </div>
          </motion.div>

          {/* Submit Buttons */}
          <motion.div
            className="gin-create-actions"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.7 }}
          >
            <motion.button
              type="button"
              onClick={() => navigate('/gins')}
              className="gin-create-btn gin-create-btn--secondary"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              Abbrechen
            </motion.button>
            <motion.button
              type="submit"
              disabled={isLoading}
              className="gin-create-btn gin-create-btn--primary"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              {isLoading ? (
                <>
                  <Loader2 className="gin-create-spinner" />
                  Speichern...
                </>
              ) : (
                <>
                  <Save />
                  Gin speichern
                </>
              )}
            </motion.button>
          </motion.div>
        </form>
      </div>
    </div>
  );
};

export default GinCreate;
