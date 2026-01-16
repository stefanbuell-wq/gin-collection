import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { adminApi } from '../api';
import {
  Server,
  HardDrive,
  Cpu,
  Clock,
  Play,
  RefreshCw,
  Terminal,
  Rocket,
  Layout,
  Shield,
  Database,
  Loader2,
  AlertTriangle,
  CheckCircle2,
  XCircle,
  ChevronDown,
  ChevronUp,
  Power
} from 'lucide-react';
import './ServerManagement.css';

interface ContainerStatus {
  name: string;
  status: string;
  health: string;
  ports: string;
}

interface ServerStatus {
  containers: ContainerStatus[];
  disk_usage: string;
  memory_usage: string;
  uptime: string;
}

interface QuickAction {
  id: string;
  name: string;
  description: string;
  icon: string;
  dangerous: boolean;
}

interface CommandResult {
  success: boolean;
  output: string;
  error?: string;
}

const iconMap: Record<string, React.ComponentType<any>> = {
  rocket: Rocket,
  server: Server,
  layout: Layout,
  shield: Shield,
  'refresh-cw': RefreshCw,
  power: Power
};

export default function ServerManagement() {
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [actions, setActions] = useState<QuickAction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [executing, setExecuting] = useState<string | null>(null);
  const [result, setResult] = useState<CommandResult | null>(null);
  const [showLogs, setShowLogs] = useState<string | null>(null);
  const [logs, setLogs] = useState<string>('');
  const [logsLoading, setLogsLoading] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      const [statusRes, actionsRes] = await Promise.all([
        adminApi.getServerStatus(),
        adminApi.getQuickActions()
      ]);
      setStatus(statusRes.data);
      setActions(actionsRes.data.actions || []);
    } catch (err: any) {
      if (err.response?.status === 404) {
        setError('Server Management ist auf diesem Server nicht aktiviert. Setze PROJECT_PATH in der .env Datei.');
      } else {
        setError('Daten konnten nicht geladen werden');
      }
    } finally {
      setLoading(false);
    }
  };

  const executeAction = async (actionId: string) => {
    if (executing) return;

    const action = actions.find(a => a.id === actionId);
    if (action?.dangerous && !confirm(`${action.name} wirklich ausführen? Dies kann den Service unterbrechen.`)) {
      return;
    }

    setExecuting(actionId);
    setResult(null);

    try {
      const response = await adminApi.executeAction(actionId);
      setResult(response.data);
      // Reload status after action
      setTimeout(loadData, 2000);
    } catch (err: any) {
      setResult({
        success: false,
        output: '',
        error: err.response?.data?.error || 'Aktion fehlgeschlagen'
      });
    } finally {
      setExecuting(null);
    }
  };

  const loadLogs = async (service: string) => {
    if (showLogs === service) {
      setShowLogs(null);
      return;
    }

    setShowLogs(service);
    setLogsLoading(true);
    setLogs('');

    try {
      const response = await adminApi.getServiceLogs(service, 100);
      setLogs(response.data.output || 'Keine Logs verfügbar');
    } catch (err) {
      setLogs('Fehler beim Laden der Logs');
    } finally {
      setLogsLoading(false);
    }
  };

  const restartService = async (service: string) => {
    if (executing) return;
    if (!confirm(`${service} wirklich neu starten?`)) return;

    setExecuting(service);
    try {
      await adminApi.restartService(service);
      setTimeout(loadData, 3000);
    } catch (err) {
      console.error('Restart failed:', err);
    } finally {
      setExecuting(null);
    }
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.08 }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { type: 'spring', stiffness: 100, damping: 15 }
    }
  };

  const getHealthIcon = (health: string) => {
    switch (health) {
      case 'healthy':
      case 'running':
        return <CheckCircle2 className="health-icon health-icon--healthy" />;
      case 'unhealthy':
        return <AlertTriangle className="health-icon health-icon--unhealthy" />;
      default:
        return <XCircle className="health-icon health-icon--stopped" />;
    }
  };

  const getServiceName = (containerName: string) => {
    const names: Record<string, string> = {
      'gin-collection-api': 'api',
      'gin-collection-frontend': 'frontend',
      'gin-collection-admin': 'admin-frontend',
      'gin-collection-mysql': 'mysql',
      'gin-collection-redis': 'redis'
    };
    return names[containerName] || containerName;
  };

  if (loading) {
    return (
      <div className="server-management">
        <div className="server-management-loader">
          <Loader2 className="server-management-loader__icon" />
          <span>Lade Server-Status...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="server-management">
        <div className="server-management-error">
          <AlertTriangle />
          <span>{error}</span>
          <button onClick={loadData} className="server-management-error__retry">
            <RefreshCw />
            Erneut versuchen
          </button>
        </div>
      </div>
    );
  }

  return (
    <motion.div
      className="server-management"
      variants={containerVariants}
      initial="hidden"
      animate="visible"
    >
      {/* Header */}
      <motion.div className="server-management-header" variants={itemVariants}>
        <div className="server-management-header__content">
          <h1 className="server-management-title">Server Management</h1>
          <p className="server-management-subtitle">
            Deployment, Services und Logs verwalten
          </p>
        </div>
        <motion.button
          className="server-management-refresh"
          onClick={loadData}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          <RefreshCw />
          Aktualisieren
        </motion.button>
      </motion.div>

      {/* System Stats */}
      <motion.div className="server-management-stats" variants={itemVariants}>
        <div className="server-management-stat">
          <HardDrive className="server-management-stat__icon" />
          <div className="server-management-stat__content">
            <span className="server-management-stat__label">Speicher</span>
            <span className="server-management-stat__value">{status?.disk_usage || 'N/A'}</span>
          </div>
        </div>
        <div className="server-management-stat">
          <Cpu className="server-management-stat__icon" />
          <div className="server-management-stat__content">
            <span className="server-management-stat__label">Memory</span>
            <span className="server-management-stat__value">{status?.memory_usage || 'N/A'}</span>
          </div>
        </div>
        <div className="server-management-stat">
          <Clock className="server-management-stat__icon" />
          <div className="server-management-stat__content">
            <span className="server-management-stat__label">Uptime</span>
            <span className="server-management-stat__value">{status?.uptime || 'N/A'}</span>
          </div>
        </div>
      </motion.div>

      {/* Quick Actions */}
      <motion.div className="server-management-section" variants={itemVariants}>
        <h2 className="server-management-section__title">
          <Rocket />
          Quick Actions
        </h2>
        <div className="server-management-actions">
          {actions.map((action) => {
            const Icon = iconMap[action.icon] || Server;
            const isExecuting = executing === action.id;

            return (
              <motion.button
                key={action.id}
                className={`server-management-action ${action.dangerous ? 'server-management-action--dangerous' : ''}`}
                onClick={() => executeAction(action.id)}
                disabled={!!executing}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <div className="server-management-action__icon">
                  {isExecuting ? <Loader2 className="spinning" /> : <Icon />}
                </div>
                <div className="server-management-action__content">
                  <span className="server-management-action__name">{action.name}</span>
                  <span className="server-management-action__desc">{action.description}</span>
                </div>
              </motion.button>
            );
          })}
        </div>
      </motion.div>

      {/* Command Result */}
      <AnimatePresence>
        {result && (
          <motion.div
            className={`server-management-result ${result.success ? 'server-management-result--success' : 'server-management-result--error'}`}
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            <div className="server-management-result__header">
              {result.success ? (
                <>
                  <CheckCircle2 />
                  <span>Erfolgreich ausgeführt</span>
                </>
              ) : (
                <>
                  <XCircle />
                  <span>Fehler: {result.error}</span>
                </>
              )}
              <button onClick={() => setResult(null)}>
                <XCircle />
              </button>
            </div>
            {result.output && (
              <pre className="server-management-result__output">{result.output}</pre>
            )}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Container Status */}
      <motion.div className="server-management-section" variants={itemVariants}>
        <h2 className="server-management-section__title">
          <Server />
          Container Status
        </h2>
        <div className="server-management-containers">
          {status?.containers?.map((container) => {
            const serviceName = getServiceName(container.name);
            const isExpanded = showLogs === serviceName;

            return (
              <motion.div
                key={container.name}
                className="server-management-container"
                layout
              >
                <div className="server-management-container__header">
                  <div className="server-management-container__info">
                    {getHealthIcon(container.health)}
                    <div className="server-management-container__details">
                      <span className="server-management-container__name">{container.name}</span>
                      <span className="server-management-container__status">{container.status}</span>
                    </div>
                  </div>
                  <div className="server-management-container__actions">
                    <motion.button
                      className="server-management-container__btn"
                      onClick={() => restartService(serviceName)}
                      disabled={!!executing}
                      title="Neu starten"
                      whileHover={{ scale: 1.1 }}
                      whileTap={{ scale: 0.9 }}
                    >
                      {executing === serviceName ? <Loader2 className="spinning" /> : <RefreshCw />}
                    </motion.button>
                    <motion.button
                      className="server-management-container__btn"
                      onClick={() => loadLogs(serviceName)}
                      title="Logs anzeigen"
                      whileHover={{ scale: 1.1 }}
                      whileTap={{ scale: 0.9 }}
                    >
                      <Terminal />
                      {isExpanded ? <ChevronUp /> : <ChevronDown />}
                    </motion.button>
                  </div>
                </div>

                <AnimatePresence>
                  {isExpanded && (
                    <motion.div
                      className="server-management-container__logs"
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      exit={{ opacity: 0, height: 0 }}
                    >
                      {logsLoading ? (
                        <div className="server-management-container__logs-loading">
                          <Loader2 className="spinning" />
                          <span>Lade Logs...</span>
                        </div>
                      ) : (
                        <pre className="server-management-container__logs-content">{logs}</pre>
                      )}
                    </motion.div>
                  )}
                </AnimatePresence>
              </motion.div>
            );
          })}
        </div>
      </motion.div>
    </motion.div>
  );
}
