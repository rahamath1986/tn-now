class CategoryModel {
  final String id;
  final String name;
  final String slug;

  const CategoryModel({
    required this.id,
    required this.name,
    required this.slug,
  });

  factory CategoryModel.fromJson(Map<String, dynamic> json) {
    return CategoryModel(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      slug: json['slug'] as String? ?? '',
    );
  }
}

class ContentCardModel {
  final String id;
  final String title;
  final String description;
  final String contentType;
  final String categoryId;
  final String districtId;
  final DateTime? publishedAt;
  final DateTime createdAt;

  const ContentCardModel({
    required this.id,
    required this.title,
    required this.description,
    required this.contentType,
    required this.categoryId,
    required this.districtId,
    this.publishedAt,
    required this.createdAt,
  });

  factory ContentCardModel.fromJson(Map<String, dynamic> json) {
    return ContentCardModel(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      description: json['description'] as String? ?? '',
      contentType: json['contentType'] as String? ?? 'PHOTO',
      categoryId: json['categoryId'] as String? ?? '',
      districtId: json['districtId'] as String? ?? '',
      publishedAt: json['publishedAt'] != null
          ? DateTime.tryParse(json['publishedAt'] as String)
          : null,
      createdAt: json['createdAt'] != null
          ? DateTime.tryParse(json['createdAt'] as String) ?? DateTime.now()
          : DateTime.now(),
    );
  }
}

class EventCardModel {
  final String id;
  final String title;
  final String description;
  final String eventName;
  final DateTime eventDate;
  final String startTime;
  final String endTime;
  final String organizerName;

  const EventCardModel({
    required this.id,
    required this.title,
    required this.description,
    required this.eventName,
    required this.eventDate,
    required this.startTime,
    required this.endTime,
    required this.organizerName,
  });

  factory EventCardModel.fromJson(Map<String, dynamic> json) {
    return EventCardModel(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      description: json['description'] as String? ?? '',
      eventName: json['eventName'] as String? ?? '',
      eventDate: json['eventDate'] != null
          ? DateTime.tryParse(json['eventDate'] as String) ?? DateTime.now()
          : DateTime.now(),
      startTime: json['startTime'] as String? ?? '',
      endTime: json['endTime'] as String? ?? '',
      organizerName: json['organizerName'] as String? ?? '',
    );
  }
}

class HomeFeedModel {
  final List<CategoryModel> categories;
  final List<ContentCardModel> trending;
  final List<ContentCardModel> latest;
  final List<EventCardModel> events;

  const HomeFeedModel({
    required this.categories,
    required this.trending,
    required this.latest,
    required this.events,
  });

  factory HomeFeedModel.fromJson(Map<String, dynamic> json) {
    final categoriesList = (json['categories'] as List? ?? [])
        .map((e) => CategoryModel.fromJson(e as Map<String, dynamic>))
        .toList();
    final trendingList = (json['trending'] as List? ?? [])
        .map((e) => ContentCardModel.fromJson(e as Map<String, dynamic>))
        .toList();
    final latestList = (json['latest'] as List? ?? [])
        .map((e) => ContentCardModel.fromJson(e as Map<String, dynamic>))
        .toList();
    final eventsList = (json['events'] as List? ?? [])
        .map((e) => EventCardModel.fromJson(e as Map<String, dynamic>))
        .toList();

    return HomeFeedModel(
      categories: categoriesList,
      trending: trendingList,
      latest: latestList,
      events: eventsList,
    );
  }

  factory HomeFeedModel.empty() {
    return const HomeFeedModel(
      categories: [],
      trending: [],
      latest: [],
      events: [],
    );
  }
}
