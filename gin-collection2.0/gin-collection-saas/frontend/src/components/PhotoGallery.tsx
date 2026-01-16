import { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Camera,
  X,
  Trash2,
  Star,
  Upload,
  Loader2,
  AlertCircle,
  Image as ImageIcon,
  Check
} from 'lucide-react';
import { photoAPI, tenantAPI } from '../api/services';
import type { GinPhoto } from '../types';
import './PhotoGallery.css';

interface PhotoGalleryProps {
  ginId: number;
  onPhotoChange?: () => void;
}

export function PhotoGallery({ ginId, onPhotoChange }: PhotoGalleryProps) {
  const [photos, setPhotos] = useState<GinPhoto[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [photoLimit, setPhotoLimit] = useState<number>(1);
  const [selectedPhoto, setSelectedPhoto] = useState<GinPhoto | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Load photos and tier limits
  useEffect(() => {
    loadPhotos();
    loadLimits();
  }, [ginId]);

  const loadPhotos = async () => {
    try {
      const response = await photoAPI.getPhotos(ginId);
      if (response.data?.data?.photos) {
        setPhotos(response.data.data.photos);
      }
    } catch (err) {
      console.error('Failed to load photos:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadLimits = async () => {
    try {
      const response = await tenantAPI.getUsage();
      if (response.data?.data?.limits?.max_photos_per_gin) {
        setPhotoLimit(response.data.data.limits.max_photos_per_gin);
      }
    } catch (err) {
      console.error('Failed to load limits:', err);
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith('image/')) {
      setError('Bitte nur Bilddateien hochladen');
      return;
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      setError('Maximale Dateigröße: 10 MB');
      return;
    }

    setUploading(true);
    setError(null);

    try {
      const response = await photoAPI.upload(ginId, file, 'bottle');
      if (response.data?.data) {
        setPhotos(prev => [...prev, response.data.data]);
        onPhotoChange?.();
      }
    } catch (err: any) {
      if (err.response?.status === 403) {
        setError(`Foto-Limit erreicht (${photos.length}/${photoLimit}). Upgrade für mehr Fotos.`);
      } else {
        setError('Fehler beim Hochladen. Bitte erneut versuchen.');
      }
      console.error('Upload failed:', err);
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleDelete = async (photoId: number) => {
    if (!confirm('Foto wirklich löschen?')) return;

    try {
      await photoAPI.delete(ginId, photoId);
      setPhotos(prev => prev.filter(p => p.id !== photoId));
      onPhotoChange?.();
    } catch (err) {
      setError('Fehler beim Löschen');
      console.error('Delete failed:', err);
    }
  };

  const handleSetPrimary = async (photoId: number) => {
    try {
      await photoAPI.setPrimary(ginId, photoId);
      setPhotos(prev => prev.map(p => ({
        ...p,
        is_primary: p.id === photoId
      })));
      onPhotoChange?.();
    } catch (err) {
      setError('Fehler beim Setzen des Hauptbildes');
      console.error('Set primary failed:', err);
    }
  };

  const canUpload = photos.length < photoLimit;

  return (
    <div className="photo-gallery">
      <div className="photo-gallery-header">
        <h3>
          <Camera size={20} />
          Fotos
        </h3>
        <div className="photo-limit-badge">
          {photos.length} / {photoLimit}
        </div>
      </div>

      {error && (
        <motion.div
          className="photo-gallery-error"
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <AlertCircle size={16} />
          <span>{error}</span>
          <button onClick={() => setError(null)}>
            <X size={14} />
          </button>
        </motion.div>
      )}

      <div className="photo-gallery-grid">
        {/* Upload Button */}
        {canUpload && (
          <motion.button
            className="photo-upload-btn"
            onClick={() => fileInputRef.current?.click()}
            disabled={uploading}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            {uploading ? (
              <>
                <Loader2 className="spinner" size={24} />
                <span>Lädt hoch...</span>
              </>
            ) : (
              <>
                <Upload size={24} />
                <span>Foto hinzufügen</span>
              </>
            )}
          </motion.button>
        )}

        {/* Photo Grid */}
        <AnimatePresence>
          {photos.map((photo) => (
            <motion.div
              key={photo.id}
              className={`photo-item ${photo.is_primary ? 'photo-item--primary' : ''}`}
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.8 }}
              layout
            >
              <img
                src={photo.photo_url}
                alt={photo.caption || 'Gin Foto'}
                onClick={() => setSelectedPhoto(photo)}
              />

              {photo.is_primary && (
                <div className="photo-primary-badge">
                  <Star size={12} />
                  Hauptbild
                </div>
              )}

              <div className="photo-actions">
                {!photo.is_primary && (
                  <button
                    className="photo-action-btn"
                    onClick={() => handleSetPrimary(photo.id)}
                    title="Als Hauptbild setzen"
                  >
                    <Star size={16} />
                  </button>
                )}
                <button
                  className="photo-action-btn photo-action-btn--delete"
                  onClick={() => handleDelete(photo.id)}
                  title="Löschen"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </motion.div>
          ))}
        </AnimatePresence>

        {/* Empty State */}
        {!loading && photos.length === 0 && !canUpload && (
          <div className="photo-empty">
            <ImageIcon size={32} />
            <p>Keine Fotos vorhanden</p>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="photo-loading">
            <Loader2 className="spinner" size={24} />
          </div>
        )}
      </div>

      {/* Limit Warning */}
      {!canUpload && photos.length > 0 && (
        <div className="photo-limit-warning">
          <AlertCircle size={16} />
          <span>Foto-Limit erreicht. <a href="/subscription">Upgrade</a> für mehr Fotos.</span>
        </div>
      )}

      {/* Hidden File Input */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        style={{ display: 'none' }}
      />

      {/* Lightbox */}
      <AnimatePresence>
        {selectedPhoto && (
          <motion.div
            className="photo-lightbox"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setSelectedPhoto(null)}
          >
            <motion.img
              src={selectedPhoto.url}
              alt={selectedPhoto.caption || 'Gin Foto'}
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.8, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
            />
            <button
              className="lightbox-close"
              onClick={() => setSelectedPhoto(null)}
            >
              <X size={24} />
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default PhotoGallery;
