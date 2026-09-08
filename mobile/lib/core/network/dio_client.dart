import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../config/app_config.dart';
import '../storage/secure_storage.dart';
import '../errors/failures.dart';

class DioClient {
  final Dio _dio;
  final SecureStorage _secureStorage;

  DioClient(this._dio, this._secureStorage) {
    _dio.options
      ..baseUrl = AppConfig.instance.apiBaseUrl
      ..connectTimeout = AppConfig.instance.connectTimeout
      ..receiveTimeout = AppConfig.instance.receiveTimeout
      ..headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      };

    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await _secureStorage.getAccessToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
        onError: (DioException e, handler) async {
          // Automatic token refresh logic placeholder (will be configured in Auth Phase)
          return handler.next(e);
        },
      ),
    );
  }

  // HTTP GET wrapper
  Future<Response> get(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.get(
        path,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _parseDioError(e);
    }
  }

  // HTTP POST wrapper
  Future<Response> post(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.post(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _parseDioError(e);
    }
  }

  // HTTP DELETE wrapper
  Future<Response> delete(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) async {
    try {
      return await _dio.delete(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
      );
    } on DioException catch (e) {
      throw _parseDioError(e);
    }
  }

  // Convert Dio exception states to clean Failure objects
  Failure _parseDioError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return const NetworkFailure('Connection timed out. Please try again.');
      case DioExceptionType.connectionError:
        return const NetworkFailure('No internet connection. Please verify your network.');
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final responseData = error.response?.data;

        if (statusCode == 401) {
          return const AuthenticationFailure('Session expired. Please log in again.');
        }

        if (statusCode == 422 && responseData is Map && responseData.containsKey('errors')) {
          final errorsMap = <String, String>{};
          final list = responseData['errors'] as List;
          for (var item in list) {
            if (item is Map && item.containsKey('field') && item.containsKey('message')) {
              errorsMap[item['field'].toString()] = item['message'].toString();
            }
          }
          return ValidationFailure(
            responseData['message']?.toString() ?? 'Validation failed',
            validationErrors: errorsMap,
          );
        }

        String errorMessage = 'Server error occurred.';
        if (responseData is Map && responseData.containsKey('message')) {
          errorMessage = responseData['message'].toString();
        }

        return ServerFailure(errorMessage, statusCode: statusCode);
      default:
        return const ServerFailure('An unexpected connection error occurred.');
    }
  }
}

// Provider mapping
final dioClientProvider = Provider<DioClient>((ref) {
  final dio = Dio();
  final secureStorage = ref.read(secureStorageProvider);
  return DioClient(dio, secureStorage);
});
