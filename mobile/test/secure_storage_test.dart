import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/storage/secure_storage.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class FakeFlutterSecureStorage extends Fake implements FlutterSecureStorage {
  final Map<String, String> _data = {};

  @override
  Future<void> write({
    required String key,
    required String? value,
    IOSOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    MacOsOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (value != null) {
      _data[key] = value;
    } else {
      _data.remove(key);
    }
  }

  @override
  Future<String?> read({
    required String key,
    IOSOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    MacOsOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    return _data[key];
  }

  @override
  Future<void> delete({
    required String key,
    IOSOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    MacOsOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    _data.remove(key);
  }

  @override
  Future<void> deleteAll({
    IOSOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    MacOsOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    _data.clear();
  }
}

void main() {
  late FakeFlutterSecureStorage fakeStorage;
  late SecureStorage secureStorage;

  setUp(() {
    fakeStorage = FakeFlutterSecureStorage();
    secureStorage = SecureStorage(fakeStorage);
  });

  test('Save and Get Access Token works', () async {
    await secureStorage.saveAccessToken('test-token');
    final token = await secureStorage.getAccessToken();
    expect(token, 'test-token');
  });

  test('Delete Access Token works', () async {
    await secureStorage.saveAccessToken('test-token');
    await secureStorage.deleteAccessToken();
    final token = await secureStorage.getAccessToken();
    expect(token, isNull);
  });

  test('Clear all works', () async {
    await secureStorage.saveAccessToken('test-token1');
    await secureStorage.saveRefreshToken('test-token2');
    await secureStorage.clearAll();
    expect(await secureStorage.getAccessToken(), isNull);
    expect(await secureStorage.getRefreshToken(), isNull);
  });
}
