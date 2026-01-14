import axios from 'axios';

export const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to attach the fake user ID for now
// In a real app, this would come from a cookie or token
api.interceptors.request.use((config) => {
    // For demo purposes, we might need a way to set this.
    // We'll rely on the server handling auth via sessions eventually,
    // but for the initial prototype we might need to be clever.
    // The current backend implementation expects X-User-ID for project listing.
    // We'll fix this later with proper Auth context.
    return config;
});
