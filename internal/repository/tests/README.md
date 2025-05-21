# Repository Integration Testing

## Testing Approaches

This directory contains two approaches to integration testing for repositories:

### 1. Traditional Integration Tests (Legacy)

Traditional integration tests connect directly to the test database and may affect each other if run concurrently.
These tests use direct database connections and may leave data in the database if a test fails unexpectedly.

Example:
```go
func TestStaffRepository_Create(t *testing.T) {
    // Connect to the test database
    testDB, err := database.NewTestDB(t)
    require.NoError(t, err)

    // Create test data
    user := createTestUser(t, testDB.DB)
    business := createTestBusiness(t, testDB.DB, user.ID)
    
    // Create the repository
    repo := repository.NewStaffRepository(testDB.DB)
    
    // Test functionality
    // ...
}
```

### 2. Transaction-Based Tests (Recommended)

Transaction-based tests provide true test isolation by running each test in a separate transaction that is automatically
rolled back when the test completes. This approach ensures tests are idempotent, can run concurrently, and don't leave
test data in the database.

Example:
```go
func TestStaffRepository_Create(t *testing.T) {
    // Create a test suite with transaction support
    suite := NewTransactionTestSuite(t)
    
    // Get repositories that use the transaction
    repos := suite.CreateTestRepositories()
    
    // Create test data
    testData := suite.CreateTestData()
    
    // Test functionality
    // ...
    
    // The transaction is automatically rolled back after the test completes
}
```

## Helper Functions

### Regular Helper Functions
- `createTestUser`: Creates a test user in the database
- `createTestBusiness`: Creates a test business in the database
- `createTestStaff`: Creates a test staff member
- `createTestStaffWithPosition`: Creates a test staff with a specific position

### Transaction-Based Helper Functions
- `createTestUserTx`: Creates a test user in a transaction
- `createTestBusinessTx`: Creates a test business in a transaction
- `createTestStaffTx`: Creates a test staff member in a transaction
- `createTestStaffWithPositionTx`: Creates a test staff with a specific position in a transaction

## Migration Guide

To migrate an existing test to use the transaction approach:

1. Replace the database connection setup with the test suite:
   ```go
   // Old approach
   testDB, err := database.NewTestDB(t)
   require.NoError(t, err)
   
   // New approach
   suite := NewTransactionTestSuite(t)
   ```

2. Replace repository creation with the test suite repositories:
   ```go
   // Old approach
   repo := repository.NewStaffRepository(testDB.DB)
   
   // New approach 
   repos := suite.CreateTestRepositories()
   staffRepo := repos.StaffRepo
   ```

3. Replace test data creation with the test suite helpers:
   ```go
   // Old approach
   user := createTestUser(t, testDB.DB)
   business := createTestBusiness(t, testDB.DB, user.ID)
   
   // New approach
   testData := suite.CreateTestData()
   user := testData.User
   business := testData.Business
   ```

## Benefits of Transaction-Based Testing

1. **True Isolation**: Each test runs in its own transaction that is rolled back after completion
2. **Idempotent Tests**: Tests can be run repeatedly with the same results
3. **Concurrency**: Tests can run concurrently without interference
4. **Performance**: Faster test execution as database cleanup is handled via rollbacks
5. **Reliability**: Tests don't leave data in the database that could affect other tests
6. **Automatic Cleanup**: No need for manual cleanup - transactions handle it automatically