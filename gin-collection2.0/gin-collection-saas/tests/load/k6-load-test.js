import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 200 }, // Ramp up to 200 users
    { duration: '5m', target: 200 }, // Stay at 200 users
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.05'],   // Error rate must be below 5%
    errors: ['rate<0.1'],             // Custom error rate below 10%
  },
};

// Base URL
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const TENANT_SUBDOMAIN = 'loadtest';
const TEST_EMAIL = 'loadtest@example.com';
const TEST_PASSWORD = 'LoadTest123!';

let authToken = '';

// Setup function - runs once before all VUs
export function setup() {
  // Register a test tenant
  const registerPayload = JSON.stringify({
    tenant_name: 'Load Test Tenant',
    subdomain: TENANT_SUBDOMAIN,
    email: TEST_EMAIL,
    password: TEST_PASSWORD,
  });

  const registerRes = http.post(`${BASE_URL}/api/v1/auth/register`, registerPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (registerRes.status === 201 || registerRes.status === 409) {
    // Login to get token
    const loginPayload = JSON.stringify({
      email: TEST_EMAIL,
      password: TEST_PASSWORD,
    });

    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, loginPayload, {
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-Subdomain': TENANT_SUBDOMAIN,
      },
    });

    if (loginRes.status === 200) {
      const data = JSON.parse(loginRes.body);
      return { token: data.token };
    }
  }

  console.error('Setup failed:', registerRes.status);
  return { token: '' };
}

// Main test function
export default function (data) {
  const token = data.token;

  if (!token) {
    console.error('No auth token available');
    errorRate.add(1);
    return;
  }

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      'X-Tenant-Subdomain': TENANT_SUBDOMAIN,
    },
  };

  // Test 1: List gins
  let res = http.get(`${BASE_URL}/api/v1/gins`, params);
  check(res, {
    'list gins status 200': (r) => r.status === 200,
    'list gins response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);

  sleep(1);

  // Test 2: Get gin stats
  res = http.get(`${BASE_URL}/api/v1/gins/stats`, params);
  check(res, {
    'stats status 200': (r) => r.status === 200,
    'stats response time < 300ms': (r) => r.timings.duration < 300,
  }) || errorRate.add(1);

  sleep(1);

  // Test 3: Create a gin
  const createPayload = JSON.stringify({
    name: `Load Test Gin ${__VU}-${__ITER}`,
    brand: 'Test Brand',
    country: 'UK',
    gin_type: 'London Dry',
    abv: 42.0,
  });

  res = http.post(`${BASE_URL}/api/v1/gins`, createPayload, params);
  check(res, {
    'create gin status 201': (r) => r.status === 201 || r.status === 403,
    'create gin response time < 500ms': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);

  if (res.status === 201) {
    const gin = JSON.parse(res.body);

    sleep(1);

    // Test 4: Get gin detail
    res = http.get(`${BASE_URL}/api/v1/gins/${gin.id}`, params);
    check(res, {
      'get gin status 200': (r) => r.status === 200,
      'get gin response time < 200ms': (r) => r.timings.duration < 200,
    }) || errorRate.add(1);

    sleep(1);

    // Test 5: Update gin
    const updatePayload = JSON.stringify({
      rating: 4,
      is_favorite: true,
    });

    res = http.put(`${BASE_URL}/api/v1/gins/${gin.id}`, updatePayload, params);
    check(res, {
      'update gin status 200': (r) => r.status === 200,
      'update gin response time < 300ms': (r) => r.timings.duration < 300,
    }) || errorRate.add(1);
  }

  sleep(2);

  // Test 6: Search gins
  res = http.get(`${BASE_URL}/api/v1/gins/search?q=Test`, params);
  check(res, {
    'search status 200': (r) => r.status === 200,
    'search response time < 400ms': (r) => r.timings.duration < 400,
  }) || errorRate.add(1);

  sleep(1);

  // Test 7: Get tenant usage
  res = http.get(`${BASE_URL}/api/v1/tenants/usage`, params);
  check(res, {
    'usage status 200': (r) => r.status === 200,
    'usage response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);

  sleep(3);
}

// Teardown function - runs once after all VUs
export function teardown(data) {
  // Optionally clean up test data
  console.log('Load test completed');
}
