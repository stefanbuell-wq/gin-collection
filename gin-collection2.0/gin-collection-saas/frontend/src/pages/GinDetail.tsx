import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useGinStore } from '../stores/ginStore';
import {
  Wine,
  ArrowLeft,
  Edit3,
  Trash2,
  Save,
  X,
  Star,
  MapPin,
  Calendar,
  Percent,
  Droplets,
  FileText,
  Camera,
  AlertTriangle,
  Check,
  Tag,
  Euro,
  Barcode,
  Clock
} from 'lucide-react';
import type { GinCreateRequest } from '../types';
import { PhotoGallery } from '../components/PhotoGallery';
import { TastingSessions } from '../components/TastingSessions';
import './GinDetail.css';

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

const COUNTRIES = [
  'Germany',
  'United Kingdom',
  'Scotland',
  'Spain',
  'Netherlands',
  'Belgium',
  'France',
  'Italy',
  'USA',
  'Japan',
  'Australia',
  'Other'
];

const BOTTLE_SIZES = [50, 100, 200, 350, 500, 700, 750, 1000, 1500];

const GinDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { currentGin, fetchGin, updateGin, deleteGin, isLoading, error, clearError } = useGinStore();

  const [isEditing, setIsEditing] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [formData, setFormData] = useState<Partial<GinCreateRequest>>({});
  const [saveSuccess, setSaveSuccess] = useState(false);

  useEffect(() => {
    if (id) {
      fetchGin(parseInt(id));
    }
  }, [id, fetchGin]);

  useEffect(() => {
    if (currentGin) {
      setFormData({
        name: currentGin.name,
        brand: currentGin.brand || '',
        country: currentGin.country || '',
        region: currentGin.region || '',
        gin_type: currentGin.gin_type || '',
        abv: currentGin.abv,
        bottle_size: currentGin.bottle_size,
        fill_level: currentGin.fill_level,
        price: currentGin.price,
        current_market_value: currentGin.current_market_value,
        barcode: currentGin.barcode || '',
        purchase_date: currentGin.purchase_date?.split('T')[0] || '',
        purchase_location: currentGin.purchase_location || '',
        rating: currentGin.rating,
        nose_notes: currentGin.nose_notes || '',
        palate_notes: currentGin.palate_notes || '',
        finish_notes: currentGin.finish_notes || '',
        general_notes: currentGin.general_notes || '',
        description: currentGin.description || '',
        recommended_tonic: currentGin.recommended_tonic || '',
        recommended_garnish: currentGin.recommended_garnish || '',
        is_finished: currentGin.is_finished
      });
    }
  }, [currentGin]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;

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

  const handleSave = async () => {
    if (!id) return;
    clearError();

    try {
      const ginData = { ...formData };
      Object.keys(ginData).forEach(key => {
        const k = key as keyof GinCreateRequest;
        if (ginData[k] === '') {
          (ginData as Record<string, unknown>)[k] = undefined;
        }
      });

      await updateGin(parseInt(id), ginData);
      setIsEditing(false);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } catch (err) {
      console.error('Failed to update gin:', err);
    }
  };

  const handleDelete = async () => {
    if (!id) return;

    try {
      await deleteGin(parseInt(id));
      navigate('/gins');
    } catch (err) {
      console.error('Failed to delete gin:', err);
    }
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    if (currentGin) {
      setFormData({
        name: currentGin.name,
        brand: currentGin.brand || '',
        country: currentGin.country || '',
        region: currentGin.region || '',
        gin_type: currentGin.gin_type || '',
        abv: currentGin.abv,
        bottle_size: currentGin.bottle_size,
        fill_level: currentGin.fill_level,
        price: currentGin.price,
        current_market_value: currentGin.current_market_value,
        barcode: currentGin.barcode || '',
        purchase_date: currentGin.purchase_date?.split('T')[0] || '',
        purchase_location: currentGin.purchase_location || '',
        rating: currentGin.rating,
        nose_notes: currentGin.nose_notes || '',
        palate_notes: currentGin.palate_notes || '',
        finish_notes: currentGin.finish_notes || '',
        general_notes: currentGin.general_notes || '',
        description: currentGin.description || '',
        recommended_tonic: currentGin.recommended_tonic || '',
        recommended_garnish: currentGin.recommended_garnish || '',
        is_finished: currentGin.is_finished
      });
    }
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 }
  };

  const renderStars = (rating: number | undefined, interactive = false) => {
    return (
      <div className="gin-detail-stars">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            onClick={() => interactive && handleRatingClick(star)}
            disabled={!interactive}
            className={`gin-detail-star ${rating && star <= rating ? 'gin-detail-star--filled' : 'gin-detail-star--empty'}`}
          >
            <Star />
          </button>
        ))}
        {rating && <span className="gin-detail-rating-text">{rating}/5</span>}
      </div>
    );
  };

  if (isLoading && !currentGin) {
    return (
      <div className="gin-detail-page">
        <div className="gin-detail-page__ambient">
          <div className="gin-detail-ambient-orb gin-detail-ambient-orb--1" />
          <div className="gin-detail-ambient-orb gin-detail-ambient-orb--2" />
        </div>
        <div className="gin-detail-container">
          <div className="gin-detail-loading">
            <div className="gin-detail-loading__spinner" />
            <span className="gin-detail-loading__text">Lade Gin-Details...</span>
          </div>
        </div>
      </div>
    );
  }

  if (!currentGin) {
    return (
      <div className="gin-detail-page">
        <div className="gin-detail-page__ambient">
          <div className="gin-detail-ambient-orb gin-detail-ambient-orb--1" />
          <div className="gin-detail-ambient-orb gin-detail-ambient-orb--2" />
        </div>
        <div className="gin-detail-container">
          <div className="gin-detail-not-found">
            <Wine />
            <h3 className="gin-detail-not-found__title">Gin nicht gefunden</h3>
            <p className="gin-detail-not-found__text">Dieser Gin existiert nicht oder wurde geloscht.</p>
            <motion.button
              onClick={() => navigate('/gins')}
              className="gin-detail-btn gin-detail-btn--primary"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              Zuruck zur Sammlung
            </motion.button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="gin-detail-page">
      {/* Ambient Background */}
      <div className="gin-detail-page__ambient">
        <div className="gin-detail-ambient-orb gin-detail-ambient-orb--1" />
        <div className="gin-detail-ambient-orb gin-detail-ambient-orb--2" />
        <div className="gin-detail-ambient-orb gin-detail-ambient-orb--3" />
      </div>

      <div className="gin-detail-container">
        {/* Header */}
        <motion.div
          className="gin-detail-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="gin-detail-header__left">
            <motion.button
              onClick={() => navigate('/gins')}
              className="gin-detail-back-btn"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <ArrowLeft size={22} />
            </motion.button>
            <div className="gin-detail-header__info">
              <h1 className="gin-detail-header__title">
                {isEditing ? 'Gin bearbeiten' : currentGin.name}
              </h1>
              {!isEditing && currentGin.brand && (
                <span className="gin-detail-header__brand">{currentGin.brand}</span>
              )}
            </div>
          </div>

          <div className="gin-detail-header__actions">
            {saveSuccess && (
              <motion.span
                className="gin-detail-success"
                initial={{ opacity: 0, x: 10 }}
                animate={{ opacity: 1, x: 0 }}
              >
                <Check size={16} />
                Gespeichert
              </motion.span>
            )}
            {isEditing ? (
              <>
                <motion.button
                  onClick={handleCancelEdit}
                  className="gin-detail-btn gin-detail-btn--secondary"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <X size={18} />
                  Abbrechen
                </motion.button>
                <motion.button
                  onClick={handleSave}
                  disabled={isLoading}
                  className="gin-detail-btn gin-detail-btn--primary"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  {isLoading ? (
                    <span className="gin-detail-spinner" />
                  ) : (
                    <Save size={18} />
                  )}
                  Speichern
                </motion.button>
              </>
            ) : (
              <>
                <motion.button
                  onClick={() => setIsEditing(true)}
                  className="gin-detail-btn gin-detail-btn--secondary"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Edit3 size={18} />
                  Bearbeiten
                </motion.button>
                <motion.button
                  onClick={() => setShowDeleteConfirm(true)}
                  className="gin-detail-btn gin-detail-btn--danger"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Trash2 size={18} />
                  Loschen
                </motion.button>
              </>
            )}
          </div>
        </motion.div>

        {/* Error */}
        {error && (
          <motion.div
            className="gin-detail-error"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
          >
            <AlertTriangle size={18} />
            <span>{error}</span>
          </motion.div>
        )}

        {/* Photo Gallery */}
        <motion.div
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          <PhotoGallery
            ginId={currentGin.id}
            onPhotoChange={() => fetchGin(currentGin.id)}
          />
        </motion.div>

        {/* Basic Information */}
        <motion.div
          className="gin-detail-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.15 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <Wine size={20} />
            </div>
            <h2 className="gin-detail-card__title">Basis-Informationen</h2>
          </div>

          {isEditing ? (
            <div className="gin-detail-form-grid">
              <div className="gin-detail-form-group gin-detail-form-group--full">
                <label className="gin-detail-form-label">Name *</label>
                <input
                  type="text"
                  name="name"
                  value={formData.name || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                  required
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Marke</label>
                <input
                  type="text"
                  name="brand"
                  value={formData.brand || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Gin-Typ</label>
                <select
                  name="gin_type"
                  value={formData.gin_type || ''}
                  onChange={handleChange}
                  className="gin-detail-select"
                >
                  <option value="">Typ wahlen...</option>
                  {GIN_TYPES.map(type => (
                    <option key={type} value={type}>{type}</option>
                  ))}
                </select>
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Land</label>
                <select
                  name="country"
                  value={formData.country || ''}
                  onChange={handleChange}
                  className="gin-detail-select"
                >
                  <option value="">Land wahlen...</option>
                  {COUNTRIES.map(country => (
                    <option key={country} value={country}>{country}</option>
                  ))}
                </select>
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Region</label>
                <input
                  type="text"
                  name="region"
                  value={formData.region || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group gin-detail-form-group--full">
                <label className="gin-detail-form-label">Beschreibung</label>
                <textarea
                  name="description"
                  value={formData.description || ''}
                  onChange={handleChange}
                  rows={2}
                  className="gin-detail-textarea"
                />
              </div>
            </div>
          ) : (
            <div className="gin-detail-grid">
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Name</span>
                <span className="gin-detail-info-value">{currentGin.name}</span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Marke</span>
                <span className={`gin-detail-info-value ${!currentGin.brand ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.brand || '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Typ</span>
                <span className={`gin-detail-info-value ${!currentGin.gin_type ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.gin_type || '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Land</span>
                <span className={`gin-detail-info-value ${!currentGin.country ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.country || '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Region</span>
                <span className={`gin-detail-info-value ${!currentGin.region ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.region || '-'}
                </span>
              </div>
              {currentGin.description && (
                <div className="gin-detail-info-item gin-detail-info-item--full">
                  <span className="gin-detail-info-label">Beschreibung</span>
                  <span className="gin-detail-info-value">{currentGin.description}</span>
                </div>
              )}
            </div>
          )}
        </motion.div>

        {/* Technical Details */}
        <motion.div
          className="gin-detail-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.2 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <Percent size={20} />
            </div>
            <h2 className="gin-detail-card__title">Technische Details</h2>
          </div>

          {isEditing ? (
            <div className="gin-detail-form-grid">
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">ABV (%)</label>
                <input
                  type="number"
                  name="abv"
                  value={formData.abv ?? ''}
                  onChange={handleChange}
                  min="0"
                  max="100"
                  step="0.1"
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Flaschengrosse (ml)</label>
                <select
                  name="bottle_size"
                  value={formData.bottle_size ?? ''}
                  onChange={handleChange}
                  className="gin-detail-select"
                >
                  <option value="">Grosse wahlen...</option>
                  {BOTTLE_SIZES.map(size => (
                    <option key={size} value={size}>{size} ml</option>
                  ))}
                </select>
              </div>
              <div className="gin-detail-form-group gin-detail-form-group--full">
                <label className="gin-detail-form-label">Fullstand (%)</label>
                <div className="gin-detail-range-wrapper">
                  <input
                    type="range"
                    name="fill_level"
                    value={formData.fill_level ?? 100}
                    onChange={handleChange}
                    min="0"
                    max="100"
                    className="gin-detail-range"
                  />
                  <span className="gin-detail-range-value">{formData.fill_level ?? 100}%</span>
                </div>
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Barcode</label>
                <input
                  type="text"
                  name="barcode"
                  value={formData.barcode || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-checkbox-label">
                  <input
                    type="checkbox"
                    name="is_finished"
                    checked={formData.is_finished || false}
                    onChange={handleChange}
                    className="gin-detail-checkbox"
                  />
                  <span className="gin-detail-checkbox-text">Flasche ist leer</span>
                </label>
              </div>
            </div>
          ) : (
            <div className="gin-detail-grid gin-detail-grid--4">
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">ABV</span>
                <span className={`gin-detail-info-value ${!currentGin.abv ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.abv ? `${currentGin.abv}%` : '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Flaschengrosse</span>
                <span className={`gin-detail-info-value ${!currentGin.bottle_size ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.bottle_size ? `${currentGin.bottle_size} ml` : '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Fullstand</span>
                <div className="gin-detail-fill-level">
                  <div className="gin-detail-fill-bar">
                    <div
                      className="gin-detail-fill-bar__progress"
                      style={{ width: `${currentGin.fill_level || 0}%` }}
                    />
                  </div>
                  <span className="gin-detail-fill-bar__text">{currentGin.fill_level || 0}%</span>
                </div>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Status</span>
                <span className={`gin-detail-info-value ${currentGin.is_finished ? 'gin-detail-info-value--muted' : 'gin-detail-info-value--success'}`}>
                  {currentGin.is_finished ? 'Leer' : 'Verfugbar'}
                </span>
              </div>
              {currentGin.barcode && (
                <div className="gin-detail-info-item">
                  <span className="gin-detail-info-label">Barcode</span>
                  <span className="gin-detail-info-value gin-detail-info-value--mono">{currentGin.barcode}</span>
                </div>
              )}
            </div>
          )}
        </motion.div>

        {/* Purchase Information */}
        <motion.div
          className="gin-detail-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.25 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <Euro size={20} />
            </div>
            <h2 className="gin-detail-card__title">Kaufinformationen</h2>
          </div>

          {isEditing ? (
            <div className="gin-detail-form-grid">
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Kaufpreis (EUR)</label>
                <input
                  type="number"
                  name="price"
                  value={formData.price ?? ''}
                  onChange={handleChange}
                  min="0"
                  step="0.01"
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Aktueller Marktwert (EUR)</label>
                <input
                  type="number"
                  name="current_market_value"
                  value={formData.current_market_value ?? ''}
                  onChange={handleChange}
                  min="0"
                  step="0.01"
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Kaufdatum</label>
                <input
                  type="date"
                  name="purchase_date"
                  value={formData.purchase_date || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Kaufort</label>
                <input
                  type="text"
                  name="purchase_location"
                  value={formData.purchase_location || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                />
              </div>
            </div>
          ) : (
            <div className="gin-detail-grid gin-detail-grid--4">
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Kaufpreis</span>
                <span className={`gin-detail-info-value ${!currentGin.price ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.price ? `${currentGin.price.toFixed(2)} EUR` : '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Marktwert</span>
                <span className={`gin-detail-info-value ${!currentGin.current_market_value ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.current_market_value ? `${currentGin.current_market_value.toFixed(2)} EUR` : '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Kaufdatum</span>
                <span className={`gin-detail-info-value ${!currentGin.purchase_date ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.purchase_date
                    ? new Date(currentGin.purchase_date).toLocaleDateString('de-DE')
                    : '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Kaufort</span>
                <span className={`gin-detail-info-value ${!currentGin.purchase_location ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.purchase_location || '-'}
                </span>
              </div>
            </div>
          )}
        </motion.div>

        {/* Rating */}
        <motion.div
          className="gin-detail-card gin-detail-card--accent"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.3 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <Star size={20} />
            </div>
            <h2 className="gin-detail-card__title">Bewertung</h2>
          </div>

          {renderStars(isEditing ? formData.rating : currentGin.rating, isEditing)}
        </motion.div>

        {/* Tasting Notes */}
        <motion.div
          className="gin-detail-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.35 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <FileText size={20} />
            </div>
            <h2 className="gin-detail-card__title">Verkostungsnotizen</h2>
          </div>

          {isEditing ? (
            <div className="gin-detail-notes">
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Nase</label>
                <textarea
                  name="nose_notes"
                  value={formData.nose_notes || ''}
                  onChange={handleChange}
                  rows={2}
                  className="gin-detail-textarea"
                  placeholder="Aromen und Dufte..."
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Gaumen</label>
                <textarea
                  name="palate_notes"
                  value={formData.palate_notes || ''}
                  onChange={handleChange}
                  rows={2}
                  className="gin-detail-textarea"
                  placeholder="Geschmack und Flavors..."
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Abgang</label>
                <textarea
                  name="finish_notes"
                  value={formData.finish_notes || ''}
                  onChange={handleChange}
                  rows={2}
                  className="gin-detail-textarea"
                  placeholder="Nachgeschmack und Finish..."
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Allgemeine Notizen</label>
                <textarea
                  name="general_notes"
                  value={formData.general_notes || ''}
                  onChange={handleChange}
                  rows={2}
                  className="gin-detail-textarea"
                  placeholder="Weitere Beobachtungen..."
                />
              </div>
            </div>
          ) : (
            <div className="gin-detail-notes">
              {currentGin.nose_notes && (
                <div className="gin-detail-note">
                  <span className="gin-detail-note__label">Nase</span>
                  <span className="gin-detail-note__value">{currentGin.nose_notes}</span>
                </div>
              )}
              {currentGin.palate_notes && (
                <div className="gin-detail-note">
                  <span className="gin-detail-note__label">Gaumen</span>
                  <span className="gin-detail-note__value">{currentGin.palate_notes}</span>
                </div>
              )}
              {currentGin.finish_notes && (
                <div className="gin-detail-note">
                  <span className="gin-detail-note__label">Abgang</span>
                  <span className="gin-detail-note__value">{currentGin.finish_notes}</span>
                </div>
              )}
              {currentGin.general_notes && (
                <div className="gin-detail-note">
                  <span className="gin-detail-note__label">Allgemeine Notizen</span>
                  <span className="gin-detail-note__value">{currentGin.general_notes}</span>
                </div>
              )}
              {!currentGin.nose_notes && !currentGin.palate_notes && !currentGin.finish_notes && !currentGin.general_notes && (
                <p className="gin-detail-note gin-detail-note--empty">Noch keine Verkostungsnotizen</p>
              )}
            </div>
          )}
        </motion.div>

        {/* Tasting Sessions */}
        <motion.div
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.38 }}
        >
          <TastingSessions ginId={currentGin.id} ginName={currentGin.name} />
        </motion.div>

        {/* Serving Suggestions */}
        <motion.div
          className="gin-detail-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.4 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon">
              <Droplets size={20} />
            </div>
            <h2 className="gin-detail-card__title">Serviervorschlage</h2>
          </div>

          {isEditing ? (
            <div className="gin-detail-form-grid">
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Empfohlenes Tonic</label>
                <input
                  type="text"
                  name="recommended_tonic"
                  value={formData.recommended_tonic || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                  placeholder="z.B. Fever-Tree Mediterranean"
                />
              </div>
              <div className="gin-detail-form-group">
                <label className="gin-detail-form-label">Empfohlene Garnitur</label>
                <input
                  type="text"
                  name="recommended_garnish"
                  value={formData.recommended_garnish || ''}
                  onChange={handleChange}
                  className="gin-detail-input"
                  placeholder="z.B. Grapefruit-Schale"
                />
              </div>
            </div>
          ) : (
            <div className="gin-detail-grid gin-detail-grid--2">
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Empfohlenes Tonic</span>
                <span className={`gin-detail-info-value ${!currentGin.recommended_tonic ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.recommended_tonic || '-'}
                </span>
              </div>
              <div className="gin-detail-info-item">
                <span className="gin-detail-info-label">Empfohlene Garnitur</span>
                <span className={`gin-detail-info-value ${!currentGin.recommended_garnish ? 'gin-detail-info-value--muted' : ''}`}>
                  {currentGin.recommended_garnish || '-'}
                </span>
              </div>
            </div>
          )}
        </motion.div>

        {/* Metadata */}
        <motion.div
          className="gin-detail-card gin-detail-card--muted"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.4, delay: 0.45 }}
        >
          <div className="gin-detail-card__header">
            <div className="gin-detail-card__icon" style={{ background: 'rgba(255,255,255,0.1)', color: 'var(--text-muted)' }}>
              <Clock size={20} />
            </div>
            <h2 className="gin-detail-card__title" style={{ color: 'var(--text-secondary)' }}>Metadaten</h2>
          </div>

          <div className="gin-detail-metadata-grid">
            <div className="gin-detail-metadata-item">
              <span className="gin-detail-metadata-label">Hinzugefugt</span>
              <span className="gin-detail-metadata-value">
                {new Date(currentGin.created_at).toLocaleDateString('de-DE', {
                  day: '2-digit',
                  month: '2-digit',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </span>
            </div>
            <div className="gin-detail-metadata-item">
              <span className="gin-detail-metadata-label">Zuletzt aktualisiert</span>
              <span className="gin-detail-metadata-value">
                {new Date(currentGin.updated_at).toLocaleDateString('de-DE', {
                  day: '2-digit',
                  month: '2-digit',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </span>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {showDeleteConfirm && (
          <motion.div
            className="gin-detail-modal-overlay"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              className="gin-detail-modal"
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
            >
              <div className="gin-detail-modal__header">
                <div className="gin-detail-modal__header-icon">
                  <AlertTriangle size={24} />
                </div>
                <h3>Gin loschen?</h3>
              </div>

              <div className="gin-detail-modal__body">
                <p>
                  Bist du sicher, dass du "{currentGin.name}" loschen mochtest?
                  Diese Aktion kann nicht ruckgangig gemacht werden.
                </p>
              </div>

              <div className="gin-detail-modal__actions">
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  className="gin-detail-modal__btn gin-detail-modal__btn--secondary"
                >
                  Abbrechen
                </button>
                <button
                  onClick={handleDelete}
                  className="gin-detail-modal__btn gin-detail-modal__btn--danger"
                >
                  Loschen
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default GinDetail;
