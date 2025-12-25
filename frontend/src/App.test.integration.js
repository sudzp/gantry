import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import App from './App';
import apiService from './services/apiService';

// Mock the API service
jest.mock('./services/apiService', () => ({
  getWorkflows: jest.fn(),
  getRuns: jest.fn(),
  getWorkflowStats: jest.fn(),
  getWorkflowRuns: jest.fn(),
  getRun: jest.fn(),
  uploadWorkflow: jest.fn(),
  triggerWorkflow: jest.fn(),
  deleteWorkflow: jest.fn(),
}));

describe('App Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    apiService.getWorkflows.mockResolvedValue([
      {
        name: 'Build',
        jobs: {
          build: {
            'runs-on': 'ubuntu',
            steps: [{ name: 'Build', run: 'echo building' }],
          },
        },
      },
    ]);
    apiService.getRuns.mockResolvedValue([
      {
        id: 'run-123',
        workflow_name: 'Build',
        status: 'success',
        started_at: new Date().toISOString(),
        jobs: {},
      },
    ]);
    apiService.getWorkflowStats.mockResolvedValue({
      total_runs: 5,
      successful_runs: 4,
      failed_runs: 1,
      success_rate: 80,
      average_duration: 120,
    });
  });

  test('renders main dashboard with workflows', async () => {
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Workflows')).toBeInTheDocument();
    });
  });

  test('displays workflow tiles with stats', async () => {
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });
  });

  test('triggers workflow on Run button click', async () => {
    apiService.triggerWorkflow.mockResolvedValue({
      id: 'run-456',
      workflow_name: 'Build',
      status: 'running',
      started_at: new Date().toISOString(),
      jobs: {},
    });

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });

    const runButton = screen.getAllByText('Run')[0];
    fireEvent.click(runButton);

    await waitFor(() => {
      expect(apiService.triggerWorkflow).toHaveBeenCalledWith('Build');
    });
  });

  test('displays workflow history on History button click', async () => {
    apiService.getWorkflowRuns.mockResolvedValue([
      {
        id: 'run-123',
        workflow_name: 'Build',
        status: 'success',
        started_at: new Date().toISOString(),
        completed_at: new Date(Date.now() + 120000).toISOString(),
        jobs: {},
      },
      {
        id: 'run-124',
        workflow_name: 'Build',
        status: 'failed',
        started_at: new Date(Date.now() - 300000).toISOString(),
        completed_at: new Date(Date.now() - 180000).toISOString(),
        jobs: {},
      },
    ]);

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });

    const historyButtons = screen.getAllByText('History');
    fireEvent.click(historyButtons[0]);

    await waitFor(() => {
      expect(screen.getByText(/Run History: Build/)).toBeInTheDocument();
    });
  });

  test('deletes workflow with confirmation', async () => {
    apiService.deleteWorkflow.mockResolvedValue({
      message: 'Workflow deleted successfully',
      name: 'Build',
    });

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });

    const deleteButton = screen.getByTitle('Delete workflow');
    fireEvent.click(deleteButton);

    // Assuming there's a confirmation dialog
    const confirmButton = await screen.findByText(/confirm|yes|delete/i);
    fireEvent.click(confirmButton);

    await waitFor(() => {
      expect(apiService.deleteWorkflow).toHaveBeenCalledWith('Build');
    });
  });

  test('navigates to run details on View Logs click', async () => {
    apiService.getRun.mockResolvedValue({
      id: 'run-123',
      workflow_name: 'Build',
      status: 'success',
      started_at: new Date().toISOString(),
      completed_at: new Date(Date.now() + 120000).toISOString(),
      jobs: {
        build: {
          status: 'success',
          steps: [
            {
              name: 'Build',
              status: 'success',
              output: 'Build successful',
              started_at: new Date().toISOString(),
              ended_at: new Date(Date.now() + 120000).toISOString(),
            },
          ],
        },
      },
    });

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });

    const logsButton = screen.getByText('Logs');
    fireEvent.click(logsButton);

    await waitFor(() => {
      expect(apiService.getRun).toHaveBeenCalled();
    });
  });

  test('uploads new workflow', async () => {
    apiService.uploadWorkflow.mockResolvedValue({
      message: 'Workflow uploaded successfully',
      name: 'Deploy',
    });

    render(<App />);
    
    const uploadInput = screen.getByRole('textbox', { hidden: true, name: /upload/i });
    
    const file = new File(['name: Deploy\njobs:\n  deploy:\n    runs-on: ubuntu\n    steps:\n      - run: echo deploy'], 'deploy.yml', { type: 'text/yaml' });
    
    fireEvent.change(uploadInput, { target: { files: [file] } });

    await waitFor(() => {
      expect(apiService.uploadWorkflow).toHaveBeenCalled();
    });
  });

  test('displays stats correctly', async () => {
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Build')).toBeInTheDocument();
    });
    
    expect(screen.getByText(/80%/)).toBeInTheDocument(); // success rate
    expect(screen.getByText(/5 runs/)).toBeInTheDocument(); // total runs
  });
});
