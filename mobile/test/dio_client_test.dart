import 'dart:convert';
import 'dart:typed_data';
import 'package:flutter_test/flutter_test.dart';
import 'package:dio/dio.dart';
import 'package:mobile/core/network/dio_client.dart';
import 'package:mobile/core/storage/secure_storage.dart';
import 'package:mobile/core/errors/failures.dart';
import 'secure_storage_test.dart'; // reuse FakeFlutterSecureStorage

class FakeHttpClientAdapter implements HttpClientAdapter {
  int statusCode = 200;
  dynamic responseData = {};
  DioException? exception;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    if (exception != null) {
      throw exception!;
    }
    
    final bytes = utf8.encode(jsonEncode(responseData));
    return ResponseBody.fromBytes(
      bytes,
      statusCode,
      headers: {
        Headers.contentTypeHeader: [Headers.jsonContentType],
      },
    );
  }

  @override
  void close({bool force = false}) {}
}

void main() {
  late Dio dio;
  late FakeHttpClientAdapter fakeAdapter;
  late SecureStorage secureStorage;
  late DioClient dioClient;

  setUp(() {
    dio = Dio();
    fakeAdapter = FakeHttpClientAdapter();
    dio.httpClientAdapter = fakeAdapter;
    secureStorage = SecureStorage(FakeFlutterSecureStorage());
    dioClient = DioClient(dio, secureStorage);
  });

  test('Successful GET request returns data', () async {
    fakeAdapter.statusCode = 200;
    fakeAdapter.responseData = {'success': true, 'data': 'hello'};

    final response = await dioClient.get('/test');
    expect(response.data['data'], 'hello');
  });

  test('401 bad response throws AuthenticationFailure', () async {
    fakeAdapter.exception = DioException(
      requestOptions: RequestOptions(path: '/test'),
      type: DioExceptionType.badResponse,
      response: Response(
        requestOptions: RequestOptions(path: '/test'),
        statusCode: 401,
        data: {'message': 'Unauthorized'},
      ),
    );

    expect(
      () => dioClient.get('/test'),
      throwsA(isA<AuthenticationFailure>()),
    );
  });

  test('422 bad response throws ValidationFailure with errors mapped', () async {
    fakeAdapter.exception = DioException(
      requestOptions: RequestOptions(path: '/test'),
      type: DioExceptionType.badResponse,
      response: Response(
        requestOptions: RequestOptions(path: '/test'),
        statusCode: 422,
        data: {
          'message': 'Validation failed',
          'errors': [
            {'field': 'email', 'message': 'Email is invalid'}
          ]
        },
      ),
    );

    try {
      await dioClient.get('/test');
      fail('Expected ValidationFailure');
    } catch (e) {
      expect(e, isA<ValidationFailure>());
      final failure = e as ValidationFailure;
      expect(failure.message, 'Validation failed');
      expect(failure.validationErrors?['email'], 'Email is invalid');
    }
  });

  test('Connection timeout throws NetworkFailure', () async {
    fakeAdapter.exception = DioException(
      requestOptions: RequestOptions(path: '/test'),
      type: DioExceptionType.connectionTimeout,
    );

    expect(
      () => dioClient.get('/test'),
      throwsA(isA<NetworkFailure>()),
    );
  });
}
