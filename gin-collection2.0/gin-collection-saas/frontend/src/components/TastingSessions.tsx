import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Wine,
  Calendar,
  Star,
  Plus,
  Edit3,
  Trash2,
  Save,
  X,
  User,
  FileText
} from 'lucide-react';
import { tastingAPI } from '../api/services';
import type { TastingSession, TastingSessionCreateRequest } from '../types';
import './TastingSessions.css';

interface TastingSessionsProps {
  ginId: number;
  ginName: string;
}

export const TastingSessions = ({ ginId, ginName }: TastingSessionsProps) => {
  const [sessions, setSessions] = useState<TastingSession[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  // Form state
  const [formData, setFormData] = useState<TastingSessionCreateRequest>({
    date: new Date().toISOString().split('T')[0],
    notes: '',
    rating: undefined
  });

  useEffect(() => {
    loadSessions();
  }, [ginId]);

  const loadSessions = async () => {
    try {
      setIsLoading(true);
      const response = await tastingAPI.getSessions(ginId);
      setSessions(response.data?.sessions || []);
      setError(null);
    } catch (err) {
      console.error('Failed to load tasting sessions:', err);
      setError('Verkostungen konnten nicht geladen werden');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSession = async () => {
    try {
      await tastingAPI.createSession(ginId, formData);
      setIsCreating(false);
      setFormData({
        date: new Date().toISOString().split('T')[0],
        notes: '',
        rating: undefined
      });
      loadSessions();
    } catch (err) {
      console.error('Failed to create tasting session:', err);
      setError('Verkostung konnte nicht erstellt werden');
    }
  };

  const handleUpdateSession = async (sessionId: number) => {
    try {
      await tastingAPI.updateSession(ginId, sessionId, formData);
      setEditingId(null);
      setFormData({
        date: new Date().toISOString().split('T')[0],
        notes: '',
        rating: undefined
      });
      loadSessions();
    } catch (err) {
      console.error('Failed to update tasting session:', err);
      setError('Verkostung konnte nicht aktualisiert werden');
    }
  };

  const handleDeleteSession = async (sessionId: number) => {
    if (!confirm('Diese Verkostung wirklich loschen?')) return;

    try {
      await tastingAPI.deleteSession(ginId, sessionId);
      loadSessions();
    } catch (err) {
      console.error('Failed to delete tasting session:', err);
      setError('Verkostung konnte nicht geloscht werden');
    }
  };

  const startEditing = (session: TastingSession) => {
    setEditingId(session.id);
    setFormData({
      date: session.date.split('T')[0],
      notes: session.notes || '',
      rating: session.rating
    });
    setIsCreating(false);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setIsCreating(false);
    setFormData({
      date: new Date().toISOString().split('T')[0],
      notes: '',
      rating: undefined
    });
  };

  const renderStars = (rating: number | undefined, interactive = false) => {
    return (
      <div className="tasting-stars">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            onClick={() => interactive && setFormData(prev => ({
              ...prev,
              rating: prev.rating === star ? undefined : star
            }))}
            disabled={!interactive}
            className={`tasting-star ${rating && star <= rating ? 'tasting-star--filled' : 'tasting-star--empty'}`}
          >
            <Star size={16} />
          </button>
        ))}
      </div>
    );
  };

  const renderForm = (isEdit = false, sessionId?: number) => (
    <motion.div
      className="tasting-form"
      initial={{ opacity: 0, height: 0 }}
      animate={{ opacity: 1, height: 'auto' }}
      exit={{ opacity: 0, height: 0 }}
    >
      <div className="tasting-form__header">
        <h4>{isEdit ? 'Verkostung bearbeiten' : 'Neue Verkostung'}</h4>
      </div>

      <div className="tasting-form__fields">
        <div className="tasting-form__field">
          <label>
            <Calendar size={14} />
            Datum
          </label>
          <input
            type="date"
            value={formData.date}
            onChange={(e) => setFormData(prev => ({ ...prev, date: e.target.value }))}
            className="tasting-input"
          />
        </div>

        <div className="tasting-form__field">
          <label>
            <Star size={14} />
            Bewertung
          </label>
          {renderStars(formData.rating, true)}
        </div>

        <div className="tasting-form__field tasting-form__field--full">
          <label>
            <FileText size={14} />
            Notizen
          </label>
          <textarea
            value={formData.notes || ''}
            onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
            placeholder="Aromen, Geschmack, Eindruck..."
            rows={3}
            className="tasting-textarea"
          />
        </div>
      </div>

      <div className="tasting-form__actions">
        <motion.button
          type="button"
          onClick={cancelEdit}
          className="tasting-btn tasting-btn--secondary"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <X size={16} />
          Abbrechen
        </motion.button>
        <motion.button
          type="button"
          onClick={() => isEdit && sessionId ? handleUpdateSession(sessionId) : handleCreateSession()}
          className="tasting-btn tasting-btn--primary"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Save size={16} />
          Speichern
        </motion.button>
      </div>
    </motion.div>
  );

  if (isLoading) {
    return (
      <div className="tasting-sessions tasting-sessions--loading">
        <div className="tasting-spinner" />
        <span>Lade Verkostungen...</span>
      </div>
    );
  }

  return (
    <div className="tasting-sessions">
      <div className="tasting-sessions__header">
        <div className="tasting-sessions__title">
          <Wine size={20} />
          <h3>Verkostungen</h3>
          <span className="tasting-count">{sessions.length}</span>
        </div>

        {!isCreating && editingId === null && (
          <motion.button
            onClick={() => setIsCreating(true)}
            className="tasting-btn tasting-btn--add"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Plus size={18} />
            Neue Verkostung
          </motion.button>
        )}
      </div>

      {error && (
        <div className="tasting-error">
          {error}
        </div>
      )}

      <AnimatePresence>
        {isCreating && renderForm(false)}
      </AnimatePresence>

      <div className="tasting-list">
        <AnimatePresence>
          {sessions.length === 0 && !isCreating ? (
            <motion.div
              className="tasting-empty"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
            >
              <Wine size={32} />
              <p>Noch keine Verkostungen</p>
              <span>Dokumentiere deine Geschmackserlebnisse</span>
            </motion.div>
          ) : (
            sessions.map((session) => (
              <motion.div
                key={session.id}
                className="tasting-item"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                layout
              >
                {editingId === session.id ? (
                  renderForm(true, session.id)
                ) : (
                  <>
                    <div className="tasting-item__header">
                      <div className="tasting-item__date">
                        <Calendar size={14} />
                        {new Date(session.date).toLocaleDateString('de-DE', {
                          day: '2-digit',
                          month: 'long',
                          year: 'numeric'
                        })}
                      </div>
                      {session.rating && (
                        <div className="tasting-item__rating">
                          {renderStars(session.rating)}
                        </div>
                      )}
                    </div>

                    {session.notes && (
                      <div className="tasting-item__notes">
                        {session.notes}
                      </div>
                    )}

                    <div className="tasting-item__footer">
                      {session.user_name && (
                        <div className="tasting-item__user">
                          <User size={12} />
                          {session.user_name}
                        </div>
                      )}

                      <div className="tasting-item__actions">
                        <motion.button
                          onClick={() => startEditing(session)}
                          className="tasting-action-btn"
                          whileHover={{ scale: 1.1 }}
                          whileTap={{ scale: 0.9 }}
                          title="Bearbeiten"
                        >
                          <Edit3 size={14} />
                        </motion.button>
                        <motion.button
                          onClick={() => handleDeleteSession(session.id)}
                          className="tasting-action-btn tasting-action-btn--danger"
                          whileHover={{ scale: 1.1 }}
                          whileTap={{ scale: 0.9 }}
                          title="Loschen"
                        >
                          <Trash2 size={14} />
                        </motion.button>
                      </div>
                    </div>
                  </>
                )}
              </motion.div>
            ))
          )}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default TastingSessions;
