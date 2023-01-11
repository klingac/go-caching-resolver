package main

func TestCachingResolver_LookupIPAddr(t *testing.T) {
	type fields struct {
		Resolver        Resolver
		Cache           Cache
		CacheExpiration int
		Cachesize       int
	}
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []net.IPAddr
		wantErr bool
	}{
		{
			name: "Test when Set raise error",
			fields: fields{
				Cache:    &MockCacheWithSetErrors{},
				Resolver: &MockResolverWithNoErrors{},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test when LookupIPAddr raise error",
			fields: fields{
				Cache:    &MockCacheWithGetErrors{},
				Resolver: &MockResolverWithLookupErrors{},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test when Get return item from cache",
			fields: fields{
				Cache:    &MockCacheWithNoErrors{},
				Resolver: &MockResolverWithNoErrors{},
			},
			args:    args{},
			want:    testIpListResolver,
			wantErr: false,
		},
		{
			name: "Test when Get return no item but Lookup will return ipAddr",
			fields: fields{
				Cache:    &MockCacheWithGetErrors{},
				Resolver: &MockResolverWithNoErrors{},
			},
			args:    args{},
			want:    testIpListResolver,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cr := &cachingResolver{
				resolver:        tt.fields.Resolver,
				cache:           tt.fields.Cache,
				cacheExpiration: tt.fields.CacheExpiration,
			}
			got, err := cr.LookupIPAddr(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupIPAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupIPAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
