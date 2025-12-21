import React from 'react';
import { CheckCircle, XCircle, Clock, Loader } from 'lucide-react';

const getStatusIcon = (status) => {
  switch (status) {
    case 'success':
      return <CheckCircle className="w-5 h-5 text-green-500" />;
    case 'failed':
      return <XCircle className="w-5 h-5 text-red-500" />;
    case 'running':
      return <Loader className="w-5 h-5 text-blue-500 animate-spin" />;
    default:
      return <Clock className="w-5 h-5 text-gray-400" />;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'success':
      return 'bg-green-100 text-green-800';
    case 'failed':
      return 'bg-red-100 text-red-800';
    case 'running':
      return 'bg-blue-100 text-blue-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

export default function RunDetails({ run }) {
  if (!run) return null;

  const renderJobs = () => {
    const jobsToRender = run.job_order && run.job_order.length > 0
      ? run.job_order.map(jobName => ({ name: jobName, job: run.jobs[jobName] })).filter(j => j.job)
      : Object.entries(run.jobs || {}).map(([name, job]) => ({ name, job }));

    return jobsToRender.map(({ name, job }) => {
      const duration = job.started_at && job.ended_at
        ? Math.round((new Date(job.ended_at) - new Date(job.started_at)) / 1000)
        : null;

      return (
        <div key={name} className="border rounded-lg p-4">
          <div className="flex items-center gap-2 mb-3">
            {getStatusIcon(job.status)}
            <h3 className="font-semibold">{name}</h3>
            <span className={`ml-auto px-2 py-1 rounded text-xs ${getStatusColor(job.status)}`}>
              {job.status}
            </span>
          </div>

          {(job.started_at || duration) && (
            <div className="text-xs text-gray-500 mb-3 space-y-1">
              {job.started_at && (
                <div>Started: {new Date(job.started_at).toLocaleString()}</div>
              )}
              {job.ended_at && (
                <div>Ended: {new Date(job.ended_at).toLocaleString()}</div>
              )}
              {duration !== null && (
                <div className="font-medium text-gray-700">Duration: {duration}s</div>
              )}
            </div>
          )}

          {job.steps && job.steps.length > 0 && (
            <div className="mb-3">
              <p className="text-sm font-medium text-gray-700 mb-2">Steps:</p>
              <ul className="space-y-1">
                {job.steps.map((step, idx) => (
                  <li key={idx} className="text-sm text-gray-600 pl-4">
                    • {step.name}
                  </li>
                ))}
              </ul>
            </div>
          )}

          {job.output && (
            <div>
              <p className="text-sm font-medium text-gray-700 mb-2">Output:</p>
              <pre className="bg-gray-900 text-green-400 p-3 rounded text-xs overflow-x-auto max-h-96">
                {job.output}
              </pre>
            </div>
          )}
        </div>
      );
    });
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold mb-4">Run Details: {run.id}</h2>

      <div className="mb-4">
        <div className="flex items-center gap-2 mb-2">
          {getStatusIcon(run.status)}
          <span className="font-medium">Status: {run.status}</span>
        </div>
        <div className="text-sm text-gray-600">
          Workflow: {run.workflow_name}
        </div>
      </div>

      <div className="space-y-4">
        {renderJobs()}
      </div>
    </div>
  );
}