import { useState, useEffect } from 'react';
import axios from 'axios';

export default function Dashboard() {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAlerts = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('http://localhost:8080/api/alerts', {
          headers: { Authorization: `Bearer ${token}` }
        });

        // Normalize the alert keys to lowercase for consistency
        const normalized = response.data.alerts.map((alert) => ({
          alertname: alert.AlertName,
          status: alert.Status,
          severity: alert.Severity,
          description: alert.Description,
          instance: alert.Instance,
          startsAt: alert.StartsAt,
          endsAt: alert.EndsAt,
        }));

        setAlerts(normalized);
      } catch (error) {
        console.error('Error fetching alerts:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAlerts();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <h2>Alerts Dashboard</h2>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ backgroundColor: '#f2f2f2' }}>
            <th style={thStyle}>Alert Name</th>
            <th style={thStyle}>Status</th>
            <th style={thStyle}>Severity</th>
            <th style={thStyle}>Description</th>
            <th style={thStyle}>Instance</th>
            <th style={thStyle}>Starts At</th>
            <th style={thStyle}>Ends At</th>
          </tr>
        </thead>
        <tbody>
          {alerts.map((alert, index) => (
            <tr key={index} style={{ border: '1px solid #ddd' }}>
              <td style={tdStyle}>{alert.alertname}</td>
              <td style={tdStyle}>{alert.status}</td>
              <td style={tdStyle}>{alert.severity}</td>
              <td style={tdStyle}>{alert.description}</td>
              <td style={tdStyle}>{alert.instance}</td>
              <td style={tdStyle}>{new Date(alert.startsAt).toLocaleString()}</td>
              <td style={tdStyle}>
                {alert.endsAt === '0001-01-01T00:00:00Z'
                  ? 'N/A'
                  : new Date(alert.endsAt).toLocaleString()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

const thStyle = {
  padding: '10px',
  border: '1px solid #ddd',
  textAlign: 'left',
};

const tdStyle = {
  padding: '10px',
  border: '1px solid #ddd',
};
