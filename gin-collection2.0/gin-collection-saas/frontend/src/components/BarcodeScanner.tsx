import { useState, useEffect, useRef } from 'react';
import { Html5Qrcode } from 'html5-qrcode';
import { motion, AnimatePresence } from 'framer-motion';
import { Camera, X, Loader2, AlertCircle, CheckCircle, Search, Wine } from 'lucide-react';
import { ginReferenceAPI } from '../api/services';
import type { GinReference } from '../types';
import './BarcodeScanner.css';

interface BarcodeScannerProps {
  onScan: (gin: GinReference) => void;
  onClose: () => void;
  isOpen: boolean;
}

interface OpenFoodFactsProduct {
  product_name?: string;
  brands?: string;
  countries?: string;
  quantity?: string;
  alcohol?: string;
  image_url?: string;
}

export function BarcodeScanner({ onScan, onClose, isOpen }: BarcodeScannerProps) {
  const [scanning, setScanning] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [scannedCode, setScannedCode] = useState<string | null>(null);
  const [searchResult, setSearchResult] = useState<{
    source: 'catalog' | 'openfoodfacts' | 'not_found';
    gin?: GinReference;
    product?: OpenFoodFactsProduct;
  } | null>(null);
  const [manualCode, setManualCode] = useState('');

  const scannerRef = useRef<Html5Qrcode | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isOpen && !scanning) {
      startScanner();
    }

    return () => {
      stopScanner();
    };
  }, [isOpen]);

  const startScanner = async () => {
    setError(null);
    setScannedCode(null);
    setSearchResult(null);

    try {
      if (scannerRef.current) {
        await stopScanner();
      }

      scannerRef.current = new Html5Qrcode('barcode-reader');

      await scannerRef.current.start(
        { facingMode: 'environment' },
        {
          fps: 10,
          qrbox: { width: 250, height: 150 },
          aspectRatio: 1.777,
        },
        onScanSuccess,
        () => {} // Ignore scan failures
      );

      setScanning(true);
    } catch (err) {
      console.error('Scanner error:', err);
      setError('Kamera konnte nicht gestartet werden. Bitte Berechtigung erteilen oder Barcode manuell eingeben.');
    }
  };

  const stopScanner = async () => {
    if (scannerRef.current) {
      try {
        await scannerRef.current.stop();
        scannerRef.current.clear();
      } catch (err) {
        // Ignore stop errors
      }
      scannerRef.current = null;
    }
    setScanning(false);
  };

  const onScanSuccess = async (decodedText: string) => {
    // Stop scanner after successful scan
    await stopScanner();
    setScannedCode(decodedText);
    await lookupBarcode(decodedText);
  };

  const lookupBarcode = async (barcode: string) => {
    setLoading(true);
    setError(null);
    setSearchResult(null);

    try {
      // Step 1: Search in our catalog
      const catalogResponse = await ginReferenceAPI.searchByBarcode(barcode);

      if (catalogResponse.data?.data) {
        setSearchResult({
          source: 'catalog',
          gin: catalogResponse.data.data,
        });
        setLoading(false);
        return;
      }
    } catch (err) {
      // Not found in catalog, continue to Open Food Facts
    }

    try {
      // Step 2: Search in Open Food Facts
      const response = await fetch(`https://world.openfoodfacts.org/api/v0/product/${barcode}.json`);
      const data = await response.json();

      if (data.status === 1 && data.product) {
        const product = data.product as OpenFoodFactsProduct;

        // Convert to GinReference format
        const gin: GinReference = {
          id: 0,
          name: product.product_name || 'Unbekannter Gin',
          brand: product.brands || undefined,
          country: product.countries?.split(',')[0] || undefined,
          abv: product.alcohol ? parseFloat(product.alcohol) : undefined,
          barcode: barcode,
          description: undefined,
        };

        setSearchResult({
          source: 'openfoodfacts',
          gin: gin,
          product: product,
        });
      } else {
        setSearchResult({ source: 'not_found' });
      }
    } catch (err) {
      console.error('Open Food Facts error:', err);
      setSearchResult({ source: 'not_found' });
    }

    setLoading(false);
  };

  const handleManualSearch = async () => {
    if (!manualCode.trim()) return;
    setScannedCode(manualCode.trim());
    await lookupBarcode(manualCode.trim());
  };

  const handleSelect = () => {
    if (searchResult?.gin) {
      onScan(searchResult.gin);
      onClose();
    }
  };

  const handleRetry = () => {
    setScannedCode(null);
    setSearchResult(null);
    setManualCode('');
    startScanner();
  };

  const handleClose = () => {
    stopScanner();
    onClose();
  };

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <motion.div
        className="barcode-scanner-overlay"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={(e) => {
          if (e.target === e.currentTarget) handleClose();
        }}
      >
        <motion.div
          className="barcode-scanner-modal"
          initial={{ opacity: 0, scale: 0.95, y: 20 }}
          animate={{ opacity: 1, scale: 1, y: 0 }}
          exit={{ opacity: 0, scale: 0.95, y: 20 }}
          transition={{ type: 'spring', damping: 25 }}
        >
          <div className="barcode-scanner-header">
            <h3>
              <Camera size={20} />
              Barcode Scanner
            </h3>
            <button className="close-btn" onClick={handleClose}>
              <X size={20} />
            </button>
          </div>

          <div className="barcode-scanner-content">
            {/* Scanner View */}
            {!scannedCode && (
              <>
                <div
                  id="barcode-reader"
                  ref={containerRef}
                  className="barcode-reader-container"
                />

                {error && (
                  <div className="scanner-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                  </div>
                )}

                <div className="manual-input-section">
                  <p className="manual-input-label">Oder Barcode manuell eingeben:</p>
                  <div className="manual-input-row">
                    <input
                      type="text"
                      value={manualCode}
                      onChange={(e) => setManualCode(e.target.value)}
                      placeholder="z.B. 4011100011298"
                      className="manual-input"
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') handleManualSearch();
                      }}
                    />
                    <button
                      onClick={handleManualSearch}
                      className="manual-search-btn"
                      disabled={!manualCode.trim()}
                    >
                      <Search size={18} />
                    </button>
                  </div>
                </div>
              </>
            )}

            {/* Loading State */}
            {loading && (
              <div className="scanner-loading">
                <Loader2 className="spinner" size={32} />
                <p>Suche nach Produkt...</p>
                <span className="barcode-display">{scannedCode}</span>
              </div>
            )}

            {/* Results */}
            {searchResult && !loading && (
              <div className="scanner-results">
                {searchResult.source === 'catalog' && searchResult.gin && (
                  <div className="result-card result-card--success">
                    <div className="result-header">
                      <CheckCircle size={24} />
                      <span>Im Katalog gefunden!</span>
                    </div>
                    <div className="result-gin">
                      <div className="result-gin-icon">
                        <Wine size={32} />
                      </div>
                      <div className="result-gin-info">
                        <h4>{searchResult.gin.name}</h4>
                        {searchResult.gin.brand && (
                          <p className="result-brand">{searchResult.gin.brand}</p>
                        )}
                        <div className="result-meta">
                          {searchResult.gin.country && <span>{searchResult.gin.country}</span>}
                          {searchResult.gin.gin_type && <span>{searchResult.gin.gin_type}</span>}
                          {searchResult.gin.abv && <span>{searchResult.gin.abv}%</span>}
                        </div>
                      </div>
                    </div>
                    <button className="select-gin-btn" onClick={handleSelect}>
                      Diesen Gin übernehmen
                    </button>
                  </div>
                )}

                {searchResult.source === 'openfoodfacts' && searchResult.gin && (
                  <div className="result-card result-card--external">
                    <div className="result-header">
                      <Search size={24} />
                      <span>Bei Open Food Facts gefunden</span>
                    </div>
                    <div className="result-gin">
                      <div className="result-gin-icon">
                        <Wine size={32} />
                      </div>
                      <div className="result-gin-info">
                        <h4>{searchResult.gin.name}</h4>
                        {searchResult.gin.brand && (
                          <p className="result-brand">{searchResult.gin.brand}</p>
                        )}
                        <div className="result-meta">
                          {searchResult.gin.country && <span>{searchResult.gin.country}</span>}
                          {searchResult.gin.abv && <span>{searchResult.gin.abv}%</span>}
                        </div>
                      </div>
                    </div>
                    <p className="external-note">
                      Daten aus externer Quelle - bitte prüfen und ggf. ergänzen
                    </p>
                    <button className="select-gin-btn" onClick={handleSelect}>
                      Daten übernehmen
                    </button>
                  </div>
                )}

                {searchResult.source === 'not_found' && (
                  <div className="result-card result-card--not-found">
                    <div className="result-header">
                      <AlertCircle size={24} />
                      <span>Produkt nicht gefunden</span>
                    </div>
                    <p className="barcode-display">{scannedCode}</p>
                    <p className="not-found-text">
                      Dieser Barcode wurde weder im Katalog noch bei Open Food Facts gefunden.
                      Du kannst die Daten manuell eingeben.
                    </p>
                  </div>
                )}

                <button className="retry-btn" onClick={handleRetry}>
                  <Camera size={18} />
                  Erneut scannen
                </button>
              </div>
            )}
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
}

export default BarcodeScanner;
