import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useGinStore } from '../stores/ginStore';
import { Wine, Star, Plus, Search } from 'lucide-react';

const GinList = () => {
  const { gins, total, fetchGins, isLoading } = useGinStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [filter, setFilter] = useState<'all' | 'available' | 'favorite'>('all');

  useEffect(() => {
    fetchGins({ filter, limit: 50 });
  }, [filter]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      fetchGins({ q: searchQuery, limit: 50 });
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">My Gins</h1>
          <p className="text-gray-600 mt-1">{total} gins in collection</p>
        </div>
        <Link to="/gins/new" className="btn btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Gin
        </Link>
      </div>

      {/* Search & Filters */}
      <div className="card">
        <div className="flex flex-col md:flex-row gap-4">
          <form onSubmit={handleSearch} className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search gins..."
                className="input pl-10"
              />
            </div>
          </form>

          <div className="flex gap-2">
            <button
              onClick={() => setFilter('all')}
              className={`btn ${filter === 'all' ? 'btn-primary' : 'btn-secondary'}`}
            >
              All
            </button>
            <button
              onClick={() => setFilter('available')}
              className={`btn ${filter === 'available' ? 'btn-primary' : 'btn-secondary'}`}
            >
              Available
            </button>
            <button
              onClick={() => setFilter('favorite')}
              className={`btn ${filter === 'favorite' ? 'btn-primary' : 'btn-secondary'}`}
            >
              Favorites
            </button>
          </div>
        </div>
      </div>

      {/* Gin Grid */}
      {isLoading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="text-gray-600 mt-4">Loading gins...</p>
        </div>
      ) : gins.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {gins.map((gin) => (
            <Link
              key={gin.id}
              to={`/gins/${gin.id}`}
              className="card hover:shadow-md transition-shadow"
            >
              {gin.primary_photo_url ? (
                <img
                  src={gin.primary_photo_url}
                  alt={gin.name}
                  className="w-full h-48 object-cover rounded-lg mb-4"
                />
              ) : (
                <div className="w-full h-48 bg-gray-200 rounded-lg mb-4 flex items-center justify-center">
                  <Wine className="w-16 h-16 text-gray-400" />
                </div>
              )}

              <h3 className="font-bold text-lg text-gray-900">{gin.name}</h3>
              <p className="text-gray-600 text-sm mt-1">
                {gin.brand} • {gin.country || 'Unknown'}
              </p>

              <div className="flex items-center justify-between mt-4">
                {gin.rating && (
                  <div className="flex items-center gap-1">
                    <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                    <span className="text-sm font-medium">{gin.rating}/5</span>
                  </div>
                )}
                {gin.price && (
                  <span className="text-sm font-medium text-gray-900">${gin.price}</span>
                )}
              </div>
            </Link>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <Wine className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No gins found</h3>
          <p className="text-gray-600 mb-4">Start building your collection</p>
          <Link to="/gins/new" className="btn btn-primary inline-flex items-center gap-2">
            <Plus className="w-5 h-5" />
            Add Your First Gin
          </Link>
        </div>
      )}
    </div>
  );
};

export default GinList;
