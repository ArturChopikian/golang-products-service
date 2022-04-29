package usecase

//func BenchmarkProductUC_Fetch(b *testing.B) {
//
//	mockProductRepos := &mocks.ProductsReposMock{}
//	mockProductRepos.On("Get", context.Background(), "Beer - Heinekin").Return(&models.Product{}, nil)
//	mockProductRepos.On("Create", context.Background(), &models.Product{}).Return(nil)
//	mockProductRepos.On("UpdatePrice", context.Background(), primitive.NewObjectID(), 1.1).Return(nil)
//
//	repos := repository.NewRepository(nil)
//	repos.Products = mockProductRepos
//
//	useCases := NewUseCases(repos)
//
//	err := useCases.ProductsUC.Fetch(context.Background(), "http://localhost:8090/list1.csv")
//	if err != nil {
//		b.Fatal(err)
//	}
//}
