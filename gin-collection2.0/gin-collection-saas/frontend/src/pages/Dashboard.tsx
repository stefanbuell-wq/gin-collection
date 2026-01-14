import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useGinStore } from '../stores/ginStore';
import { Wine, Star, DollarSign, TrendingUp, Plus } from 'lucide-react';

const Dashboard = () => {
  const { stats, fetchStats, gins, fetchGins, isLoading } = useGinStore();
  const [recentGins, setRecentGins] = useState<typeof gins>([]);

  useEffect(() => {
    fetchStats();
    fetchGins({ limit: 5, sort: 'created_at' });
  }, []);

  useEffect(() => {
    setRecentGins(gins.slice(0, 5));
  }, [gins]);

  const statCards = [
    {
      label: 'Total Gins',
      value: stats?.total_gins || 0,
      icon: Wine,
      color: 'bg-blue-500',
    },
    {
      label: 'Available',
      value: stats?.available_gins || 0,
      icon: TrendingUp,
      color: 'bg-green-500',
    },
    {
      label: 'Favorites',
      value: stats?.favorite_count || 0,
      icon: Star,
      color: 'bg-yellow-500',
    },
    {
      label: 'Total Value',
      value: `$${stats?.total_value?.toFixed(2) || '0.00'}`,
      icon: DollarSign,
      color: 'bg-purple-500',
    },
  ];

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-600 mt-1">Overview of your gin collection</p>
        </div>
        <Link to="/gins/new" className="btn btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Gin
        </Link>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="card">
              <div className="flex items-center gap-4">
                <div className={`${stat.color} p-3 rounded-lg`}>
                  <Icon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <p className="text-sm text-gray-600">{stat.label}</p>
                  <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* By Country */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">By Country</h3>
          {stats?.by_country && Object.keys(stats.by_country).length > 0 ? (
            <div className="space-y-3">
              {Object.entries(stats.by_country)
                .sort(([, a], [, b]) => b - a)
                .slice(0, 5)
                .map(([country, count]) => (
                  <div key={country} className="flex items-center justify-between">
                    <span className="text-gray-700">{country || 'Unknown'}</span>
                    <div className="flex items-center gap-3">
                      <div className="w-32 h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-primary-600 rounded-full"
                          style={{
                            width: `${(count / (stats.total_gins || 1)) * 100}%`,
                          }}
                        ></div>
                      </div>
                      <span className="text-sm font-medium text-gray-900 w-8 text-right">
                        {count}
                      </span>
                    </div>
                  </div>
                ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">No data available</p>
          )}
        </div>

        {/* By Type */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">By Type</h3>
          {stats?.by_type && Object.keys(stats.by_type).length > 0 ? (
            <div className="space-y-3">
              {Object.entries(stats.by_type)
                .sort(([, a], [, b]) => b - a)
                .slice(0, 5)
                .map(([type, count]) => (
                  <div key={type} className="flex items-center justify-between">
                    <span className="text-gray-700">{type || 'Unknown'}</span>
                    <div className="flex items-center gap-3">
                      <div className="w-32 h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-green-600 rounded-full"
                          style={{
                            width: `${(count / (stats.total_gins || 1)) * 100}%`,
                          }}
                        ></div>
                      </div>
                      <span className="text-sm font-medium text-gray-900 w-8 text-right">
                        {count}
                      </span>
                    </div>
                  </div>
                ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">No data available</p>
          )}
        </div>
      </div>

      {/* Recent Gins */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">Recently Added</h3>
          <Link to="/gins" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            View All
          </Link>
        </div>

        {isLoading ? (
          <p className="text-center py-8 text-gray-500">Loading...</p>
        ) : recentGins.length > 0 ? (
          <div className="space-y-3">
            {recentGins.map((gin) => (
              <Link
                key={gin.id}
                to={`/gins/${gin.id}`}
                className="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 transition-colors"
              >
                {gin.primary_photo_url ? (
                  <img
                    src={gin.primary_photo_url}
                    alt={gin.name}
                    className="w-12 h-12 rounded-lg object-cover"
                  />
                ) : (
                  <div className="w-12 h-12 rounded-lg bg-gray-200 flex items-center justify-center">
                    <Wine className="w-6 h-6 text-gray-400" />
                  </div>
                )}
                <div className="flex-1">
                  <h4 className="font-medium text-gray-900">{gin.name}</h4>
                  <p className="text-sm text-gray-600">
                    {gin.brand} • {gin.country || 'Unknown'}
                  </p>
                </div>
                {gin.rating && (
                  <div className="flex items-center gap-1">
                    <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                    <span className="text-sm font-medium">{gin.rating}</span>
                  </div>
                )}
              </Link>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <Wine className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-600">No gins yet</p>
            <Link to="/gins/new" className="text-primary-600 hover:text-primary-700 text-sm font-medium mt-2 inline-block">
              Add your first gin
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
