package postgres

//func Test_DatabaseRepository_Save_RepositoryError(t *testing.T) {
//	mockDb := new(storagemocks.DBIface)
//	repo := NewDatabaseRepository(mockDb, 5*time.Second, "currency_rates")
//
//	mockDb.On("Begin", mock.Anything).Return(nil, nil)
//
//	uuidStr, err := uuid.NewUUID()
//	require.NoError(t, err)
//	data := bole.CurrencyData{
//		ID: uuidStr.String(),
//		Data: map[string]bole.Currency{
//			"USD": {Code: "USD", Value: 1.0},
//			"EUR": {Code: "EUR", Value: 0.9},
//		},
//	}
//	data.Meta.LastUpdatedAt = time.Now()
//	err = repo.Save(context.Background(), &data)
//	mockDb.AssertExpectations(t)
//	assert.NoError(t, err)
//
//}

//func Test_DatabaseRepository_Save_Succeed(t *testing.T) {
//	mockPool, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockPool.Close()
//
//	mockPool.ExpectBegin()
//	mockPool.CopyFrom(pgx.Identifier{"currency_rates"}, []string{"code", "value", "last_updated_at", "batch_id"}, pgx.CopyFromRows([][]interface{}{})
//	mockPool.ExpectExec("INSERT INTO product_viewers").
//		WithArgs(2, 3).
//		WillReturnResult(pgxmock.NewResult("INSERT", 1))
//	mockPool.ExpectCommit()
//
//	// now we execute our method
//	if err = recordStats(mockPool, 2, 3); err != nil {
//		t.Errorf("error was not expected while updating: %s", err)
//	}
//
//	// we make sure that all expectations were met
//	if err := mockPool.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
//
//func Test_DatabaseRepository_SaveCall_Succeed(t *testing.T) {
//	mockPool, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	repo := NewDatabaseRepository(mockPool, 5*time.Second, "api_calls")
//
//	ctx := context.Background()
//	data := bole.ApiCall{
//		ID:           "test-call-id",
//		StatusCode:   200,
//		ResponseTime: 150,
//		Timeout:      false,
//		ErrorMessage: "",
//	}
//
//	mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgx.CommandTag{}, nil)
//
//	// Ejecutando el método bajo prueba
//	err := repo.SaveCall(ctx, data)
//
//	// Verificando los resultados
//	assert.NoError(t, err)
//	mockDB.AssertExpectations(t)
//}
