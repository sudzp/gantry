// API Service - Centralized API calls

const API_URL = process.env.REACT_APP_API_URL || "http://localhost:8080/api";

class ApiService {
  // Workflows
  async getWorkflows() {
    const response = await fetch(`${API_URL}/workflows`);
    if (!response.ok) throw new Error("Failed to fetch workflows");
    return response.json();
  }

  async uploadWorkflow(file) {
    const text = await file.text();
    const response = await fetch(`${API_URL}/workflows`, {
      method: "POST",
      body: text,
      headers: { "Content-Type": "text/yaml" },
    });
    if (!response.ok) throw new Error("Failed to upload workflow");
    return response.json();
  }

  async triggerWorkflow(name) {
    const response = await fetch(`${API_URL}/workflows/${name}/trigger`, {
      method: "POST",
    });
    if (!response.ok) throw new Error("Failed to trigger workflow");
    return response.json();
  }

  // Runs
  async getRuns() {
    const response = await fetch(`${API_URL}/runs`);
    if (!response.ok) throw new Error("Failed to fetch runs");
    const data = await response.json();
    return (data || []).sort(
      (a, b) => new Date(b.started_at) - new Date(a.started_at),
    );
  }

  async getRun(id) {
    const response = await fetch(`${API_URL}/runs/${id}`);
    if (!response.ok) throw new Error("Failed to fetch run");
    return response.json();
  }
}

// Create instance and export
const apiService = new ApiService();
export default apiService;
