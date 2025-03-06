import React, { useEffect, useState } from 'react';

function App() {
  const [connected, setConnected] = useState(false);
  const [lastScan, setLastScan] = useState('Never');
  const [wallets, setWallets] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Check API health
    fetch('/api/v1/health')
      .then(response => {
        if (response.ok) {
          setConnected(true);
          return fetch('/api/v1/wallets');
        } else {
          throw new Error('API connection failed');
        }
      })
      .then(response => response.json())
      .then(data => {
        setWallets(data);
        return fetch('/api/v1/alerts');
      })
      .then(response => response.json())
      .then(data => {
        setAlerts(data);
        return fetch('/api/v1/status');
      })
      .then(response => response.json())
      .then(data => {
        setLastScan(data.last_scan || 'Never');
        setLoading(false);
      })
      .catch(err => {
        console.error('Error:', err);
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
        <div className="container">
          <a className="navbar-brand" href="/">Insider Monitor</a>
          <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarNav">
            <ul className="navbar-nav ms-auto">
              <li className="nav-item">
                <a className="nav-link active" href="/">Dashboard</a>
              </li>
              <li className="nav-item">
                <a className="nav-link" href="/wallets">Wallets</a>
              </li>
              <li className="nav-item">
                <a className="nav-link" href="/config">Configuration</a>
              </li>
            </ul>
          </div>
        </div>
      </nav>

      <div className="container mt-4">
        <div className="row">
          <div className="col-12">
            <div className="card">
              <div className="card-header d-flex justify-content-between align-items-center">
                <h2>Dashboard</h2>
                <div>
                  <span className={`badge ${connected ? 'bg-success' : 'bg-danger'}`}>
                    {connected ? 'Connected' : 'Disconnected'}
                  </span>
                  <span className="ms-2">Last scan: {lastScan}</span>
                </div>
              </div>
              <div className="card-body">
                {error && (
                  <div className="alert alert-danger">
                    Error: {error}
                  </div>
                )}

                <h3>Monitored Wallets</h3>
                {loading ? (
                  <p>Loading wallet data...</p>
                ) : (
                  <div className="table-responsive">
                    <table className="table table-striped table-hover">
                      <thead>
                        <tr>
                          <th>Address</th>
                          <th>Label</th>
                          <th>Last Scanned</th>
                          <th>Token Count</th>
                          <th>Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {wallets.length === 0 ? (
                          <tr>
                            <td colSpan="5" className="text-center">No wallets configured</td>
                          </tr>
                        ) : (
                          wallets.map((wallet, index) => (
                            <tr key={index}>
                              <td>{wallet.address}</td>
                              <td>{wallet.label || 'Unnamed'}</td>
                              <td>{wallet.last_scanned || 'Never'}</td>
                              <td>{wallet.token_count || 0}</td>
                              <td>
                                <a href={`/wallet/${wallet.address}`} className="btn btn-sm btn-primary">View</a>
                              </td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                  </div>
                )}

                <h3>Recent Alerts</h3>
                {loading ? (
                  <p>Loading alert data...</p>
                ) : (
                  <div className="table-responsive">
                    <table className="table table-striped table-hover">
                      <thead>
                        <tr>
                          <th>Time</th>
                          <th>Wallet</th>
                          <th>Type</th>
                          <th>Message</th>
                          <th>Level</th>
                        </tr>
                      </thead>
                      <tbody>
                        {alerts.length === 0 ? (
                          <tr>
                            <td colSpan="5" className="text-center">No recent alerts</td>
                          </tr>
                        ) : (
                          alerts.map((alert, index) => {
                            const levelClass = alert.level === 'CRITICAL' ? 'danger' :
                                            (alert.level === 'WARNING' ? 'warning' : 'info');

                            return (
                              <tr key={index}>
                                <td>{alert.timestamp}</td>
                                <td>{alert.wallet_address}</td>
                                <td>{alert.alert_type}</td>
                                <td>{alert.message}</td>
                                <td>
                                  <span className={`badge bg-${levelClass}`}>{alert.level}</span>
                                </td>
                              </tr>
                            );
                          })
                        )}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
