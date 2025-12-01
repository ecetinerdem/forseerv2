import { useState, useEffect } from 'react';
import { Eye, EyeOff, Plus, Trash2, Edit2, Search, LogOut } from 'lucide-react';

const API_URL = 'http://localhost:8080/v1';

/* ========== TYPE DEFINITIONS (mirror Go backend) ========== */
interface Stock {
  id: number;
  symbol: string;
  shares: number;
  average_price: number;
  portfolio_id: number;
  created_at?: string;
  updated_at?: string;
}

interface Portfolio {
  id: number;
  user_id: number;
  name: string;
  stocks: Stock[];
  created_at?: string;
  updated_at?: string;
}

interface User {
  id: number;
  username: string;
  email: string;
  created_at?: string;
  updated_at?: string;
}

/* ================== MAIN APP COMPONENT ================== */
export default function App() {
  const [view, setView] = useState<'login' | 'register' | 'portfolios' | 'portfolio-detail'>('login');
  const [token, setToken] = useState('');
  const [user, setUser] = useState<User | null>(null);
  const [portfolios, setPortfolios] = useState<Portfolio[]>([]);
  const [selectedPortfolio, setSelectedPortfolio] = useState<Portfolio | null>(null);
  const [error, setError] = useState('');

  // Load token from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('token');
    if (saved) setToken(saved);
  }, []);

  // Fetch user + portfolios when token changes
  useEffect(() => {
    if (token) {
      fetchUser();
      fetchPortfolios();
    }
  }, [token]);

  const fetchUser = async () => {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const res = await fetch(`${API_URL}/users/${payload.sub}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data: User = await res.json();
        setUser(data);
        setView('portfolios');
      } else {
        logout();
      }
    } catch (e) {
      console.error(e);
      logout();
    }
  };

  const fetchPortfolios = async () => {
    try {
      const res = await fetch(`${API_URL}/portfolios`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data: Portfolio[] = await res.json();
        setPortfolios(data);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const logout = () => {
    setToken('');
    setUser(null);
    setPortfolios([]);
    setSelectedPortfolio(null);
    localStorage.removeItem('token');
    setView('login');
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <nav className="bg-gray-800 border-b border-gray-700 px-6 py-4">
        <div className="max-w-6xl mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold text-cyan-400">Forseer</h1>
          {user && (
            <div className="flex items-center gap-4">
              <span className="text-gray-300">{user.username}</span>
              <button onClick={logout} className="text-red-400 hover:text-red-300">
                <LogOut size={20} />
              </button>
            </div>
          )}
        </div>
      </nav>

      <div className="max-w-6xl mx-auto px-6 py-8">
        {view === 'login' && (
          <LoginView setToken={setToken} setView={setView} setError={setError} error={error} />
        )}
        {view === 'register' && (
          <RegisterView setView={setView} setError={setError} error={error} />
        )}
        {view === 'portfolios' && (
          <PortfoliosView
            portfolios={portfolios}
            
            setSelectedPortfolio={setSelectedPortfolio}
            setView={setView}
            token={token}
            fetchPortfolios={fetchPortfolios}
          />
        )}
        {view === 'portfolio-detail' && selectedPortfolio && (
          <PortfolioDetail
            portfolio={selectedPortfolio}
            setView={setView}
            token={token}
            fetchPortfolios={fetchPortfolios}
            onPortfolioUpdate={(updated) => {
              setSelectedPortfolio(updated);
              setPortfolios((prev) => prev.map((p) => (p.id === updated.id ? updated : p)));
            }}
          />
        )}
      </div>
    </div>
  );
}

/* ================== LOGIN VIEW ================== */
function LoginView({
  setToken,
  setView,
  setError,
  error,
}: {
  setToken: (t: string) => void;
  setView: (v: any) => void;
  setError: (e: string) => void;
  error: string;
}) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleLogin = async () => {
    setError('');
    try {
      const res = await fetch(`${API_URL}/authentication/token`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: '', email, password }),
      });
      if (res.ok) {
        const { token } = await res.json();
        localStorage.setItem('token', token);
        setToken(token);
      } else {
        const err = await res.text();
        setError(err || 'Invalid credentials');
      }
    } catch {
      setError('Login failed. Please try again.');
    }
  };

  return (
    <div className="max-w-md mx-auto mt-20">
      <div className="bg-gray-800 rounded-lg p-8 border border-gray-700">
        <h2 className="text-2xl font-bold mb-6 text-cyan-400">Login</h2>
        {error && <div className="bg-red-900/50 text-red-300 p-3 rounded mb-4">{error}</div>}
        <div className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleLogin()}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
          />
          <div className="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleLogin()}
              className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-2.5 text-gray-400"
            >
              {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
            </button>
          </div>
          <button
            onClick={handleLogin}
            className="w-full bg-cyan-500 hover:bg-cyan-600 text-white py-2 rounded font-medium"
          >
            Login
          </button>
        </div>
        <button
          onClick={() => setView('register')}
          className="w-full mt-4 text-cyan-400 hover:text-cyan-300"
        >
          Need an account? Register
        </button>
      </div>
    </div>
  );
}

/* ================== REGISTER VIEW ================== */
function RegisterView({
  setView,
  setError,
  error,
}: {
  setView: (v: any) => void;
  setError: (e: string) => void;
  error: string;
}) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [success, setSuccess] = useState('');

  const handleRegister = async () => {
    setError('');
    setSuccess('');
    try {
      const res = await fetch(`${API_URL}/authentication/user`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password }),
      });
      if (res.ok) {
        setSuccess('Registration successful! Check your email to activate your account.');
        setTimeout(() => setView('login'), 3000);
      } else {
        const data = await res.json();
        setError(data.error || 'Registration failed');
      }
    } catch {
      setError('Registration failed');
    }
  };

  return (
    <div className="max-w-md mx-auto mt-20">
      <div className="bg-gray-800 rounded-lg p-8 border border-gray-700">
        <h2 className="text-2xl font-bold mb-6 text-cyan-400">Register</h2>
        {error && <div className="bg-red-900/50 text-red-300 p-3 rounded mb-4">{error}</div>}
        {success && <div className="bg-green-900/50 text-green-300 p-3 rounded mb-4">{success}</div>}
        <div className="space-y-4">
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            maxLength={100}
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            maxLength={255}
          />
          <div className="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              placeholder="Password (8-16 chars)"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
              minLength={8}
              maxLength={16}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-2.5 text-gray-400"
            >
              {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
            </button>
          </div>
          <button
            onClick={handleRegister}
            className="w-full bg-cyan-500 hover:bg-cyan-600 text-white py-2 rounded font-medium"
          >
            Register
          </button>
        </div>
        <button
          onClick={() => setView('login')}
          className="w-full mt-4 text-cyan-400 hover:text-cyan-300"
        >
          Already have an account? Login
        </button>
      </div>
    </div>
  );
}

/* ================== PORTFOLIOS VIEW ================== */
function PortfoliosView({
  portfolios,
  setSelectedPortfolio,
  setView,
  token,
  fetchPortfolios,
}: {
  portfolios: Portfolio[];
  
  setSelectedPortfolio: (p: Portfolio) => void;
  setView: (v: any) => void;
  token: string;
  fetchPortfolios: () => Promise<void>;
}) {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState<Portfolio[]>([]);

  const handleSearch = async () => {
    if (!searchTerm.trim()) return;
    try {
      const res = await fetch(`${API_URL}/portfolios/search?name=${encodeURIComponent(searchTerm)}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data: Portfolio[] = await res.json();
        setSearchResults(data);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const displayed = searchTerm ? searchResults : portfolios;

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h2 className="text-3xl font-bold text-cyan-400">My Portfolios</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          className="bg-cyan-500 hover:bg-cyan-600 px-4 py-2 rounded flex items-center gap-2"
        >
          <Plus size={20} /> Create Portfolio
        </button>
      </div>

      <div className="mb-6 flex gap-2">
        <input
          type="text"
          placeholder="Search portfolios..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
          className="flex-1 bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
        />
        <button onClick={handleSearch} className="bg-cyan-500 hover:bg-cyan-600 px-4 py-2 rounded">
          <Search size={20} />
        </button>
        {searchTerm && (
          <button
            onClick={() => {
              setSearchTerm('');
              setSearchResults([]);
            }}
            className="text-gray-400 hover:text-white px-4"
          >
            Clear
          </button>
        )}
      </div>

      <div className="grid gap-4">
        {displayed.map((portfolio) => (
          <div
            key={portfolio.id}
            onClick={() => {
              setSelectedPortfolio(portfolio);
              setView('portfolio-detail');
            }}
            className="bg-gray-800 border border-gray-700 rounded-lg p-6 hover:border-cyan-400 cursor-pointer transition-colors"
          >
            <h3 className="text-xl font-semibold text-cyan-400 mb-2">{portfolio.name}</h3>
            <p className="text-gray-400">Stocks: {portfolio.stocks.length}</p>
          </div>
        ))}
        {displayed.length === 0 && (
          <p className="text-gray-400 text-center py-8">No portfolios found</p>
        )}
      </div>

      {showCreateModal && (
        <CreatePortfolioModal
          onClose={() => setShowCreateModal(false)}
          token={token}
          fetchPortfolios={fetchPortfolios}
        />
      )}
    </div>
  );
}

/* ================== CREATE PORTFOLIO MODAL ================== */
function CreatePortfolioModal({
  onClose,
  token,
  fetchPortfolios,
}: {
  onClose: () => void;
  token: string;
  fetchPortfolios: () => Promise<void>;
}) {
  const [name, setName] = useState('');

  const handleCreate = async () => {
    if (!name.trim()) return;
    try {
      const res = await fetch(`${API_URL}/portfolios`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name, stocks: [] }),
      });
      if (res.ok) {
        await fetchPortfolios();
        onClose();
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4 border border-gray-700"
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className="text-2xl font-bold mb-4 text-cyan-400">Create Portfolio</h3>
        <input
          type="text"
          placeholder="Portfolio name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleCreate()}
          className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 mb-4 focus:outline-none focus:border-cyan-400"
          maxLength={50}
        />
        <div className="flex gap-2">
          <button onClick={handleCreate} className="flex-1 bg-cyan-500 hover:bg-cyan-600 py-2 rounded">
            Create
          </button>
          <button onClick={onClose} className="flex-1 bg-gray-700 hover:bg-gray-600 py-2 rounded">
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}

/* ================== PORTFOLIO DETAIL ================== */
function PortfolioDetail({
  portfolio,
  setView,
  token,
  fetchPortfolios,
  onPortfolioUpdate,
}: {
  portfolio: Portfolio;
  setView: (v: any) => void;
  token: string;
  fetchPortfolios: () => Promise<void>;
  onPortfolioUpdate: (p: Portfolio) => void;
}) {
  const [stocks, setStocks] = useState<Stock[]>(portfolio.stocks);
  const [showAddStock, setShowAddStock] = useState(false);
  const [editingName, setEditingName] = useState(false);
  const [newName, setNewName] = useState(portfolio.name);

  const handleUpdateName = async () => {
    if (!newName.trim()) return;
    try {
      const res = await fetch(`${API_URL}/portfolios/${portfolio.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: newName }),
      });
      if (res.ok) {
        const updated: Portfolio = await res.json();
        onPortfolioUpdate(updated);
        setEditingName(false);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeletePortfolio = async () => {
    if (!confirm('Delete this portfolio?')) return;
    try {
      const res = await fetch(`${API_URL}/portfolios/${portfolio.id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        await fetchPortfolios();
        setView('portfolios');
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeleteStock = async (symbol: string) => {
    if (!confirm(`Remove ${symbol}?`)) return;
    try {
      const res = await fetch(`${API_URL}/portfolios/${portfolio.id}/stocks/${symbol}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        setStocks((prev) => prev.filter((s) => s.symbol !== symbol));
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div>
      <button onClick={() => setView('portfolios')} className="text-cyan-400 hover:text-cyan-300 mb-6">
        ← Back to Portfolios
      </button>

      <div className="bg-gray-800 border border-gray-700 rounded-lg p-6 mb-6">
        <div className="flex justify-between items-center mb-4">
          {editingName ? (
            <div className="flex gap-2 flex-1">
              <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleUpdateName()}
                className="flex-1 bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
                maxLength={50}
              />
              <button onClick={handleUpdateName} className="bg-green-500 hover:bg-green-600 px-4 py-2 rounded">
                Save
              </button>
              <button onClick={() => setEditingName(false)} className="bg-gray-700 hover:bg-gray-600 px-4 py-2 rounded">
                Cancel
              </button>
            </div>
          ) : (
            <>
              <h2 className="text-3xl font-bold text-cyan-400">{portfolio.name}</h2>
              <div className="flex gap-2">
                <button onClick={() => setEditingName(true)} className="text-cyan-400 hover:text-cyan-300">
                  <Edit2 size={20} />
                </button>
                <button onClick={handleDeletePortfolio} className="text-red-400 hover:text-red-300">
                  <Trash2 size={20} />
                </button>
              </div>
            </>
          )}
        </div>

        <button
          onClick={() => setShowAddStock(true)}
          className="bg-green-500 hover:bg-green-600 px-4 py-2 rounded flex items-center gap-2"
        >
          <Plus size={20} /> Add Stock
        </button>
      </div>

      <div className="grid gap-4">
        {stocks.map((stock) => (
          <div
            key={stock.id}
            className="bg-gray-800 border border-gray-700 rounded-lg p-6 flex justify-between items-center"
          >
            <div>
              <h3 className="text-xl font-bold text-cyan-400">{stock.symbol}</h3>
              <p className="text-gray-300">Shares: {stock.shares}</p>
              <p className="text-gray-300">Avg Price: ${stock.average_price}</p>
            </div>
            <button onClick={() => handleDeleteStock(stock.symbol)} className="text-red-400 hover:text-red-300">
              <Trash2 size={20} />
            </button>
          </div>
        ))}
        {stocks.length === 0 && (
          <p className="text-gray-400 text-center py-8">No stocks in this portfolio</p>
        )}
      </div>

      {showAddStock && (
        <AddStockModal
          portfolioId={portfolio.id}
          onClose={() => setShowAddStock(false)}
          token={token}
          onStockAdded={(newStock) => setStocks((prev) => [...prev, newStock])}
        />
      )}
    </div>
  );
}

/* ================== ADD STOCK MODAL ================== */
function AddStockModal({
  portfolioId,
  onClose,
  token,
  onStockAdded,
}: {
  portfolioId: number;
  onClose: () => void;
  token: string;
  onStockAdded: (s: Stock) => void;
}) {
  const [symbol, setSymbol] = useState('');
  const [shares, setShares] = useState('');
  const [avgPrice, setAvgPrice] = useState('');

  const handleAdd = async () => {
    if (!symbol || !shares || !avgPrice) return;
    try {
      const res = await fetch(`${API_URL}/portfolios/${portfolioId}/stocks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          symbol: symbol.toUpperCase(),
          shares: parseFloat(shares),
          average_price: parseFloat(avgPrice),
        }),
      });
      if (res.ok) {
        const newStock: Stock = await res.json();
        onStockAdded(newStock);
        onClose();
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4 border border-gray-700"
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className="text-2xl font-bold mb-4 text-cyan-400">Add Stock</h3>
        <div className="space-y-4">
          <input
            type="text"
            placeholder="Symbol (e.g. AAPL)"
            value={symbol}
            onChange={(e) => setSymbol(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            maxLength={4}
          />
          <input
            type="number"
            placeholder="Shares"
            value={shares}
            onChange={(e) => setShares(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            step="0.01"
            min="0.01"
          />
          <input
            type="number"
            placeholder="Average Price"
            value={avgPrice}
            onChange={(e) => setAvgPrice(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-2 focus:outline-none focus:border-cyan-400"
            step="0.01"
            min="0.01"
          />
          <div className="flex gap-2">
            <button onClick={handleAdd} className="flex-1 bg-green-500 hover:bg-green-600 py-2 rounded">
              Add
            </button>
            <button onClick={onClose} className="flex-1 bg-gray-700 hover:bg-gray-600 py-2 rounded">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}