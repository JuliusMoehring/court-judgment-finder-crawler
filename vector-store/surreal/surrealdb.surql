--- ANALYZER

DEFINE ANALYZER blank_snowball_ger TOKENIZERS BLANK FILTERS SNOWBALL(GERMAN);


--- PAGE

DEFINE TABLE page TYPE ANY SCHEMAFULL
	PERMISSIONS NONE
;
DEFINE FIELD page ON page TYPE int ASSERT $value > 0
	PERMISSIONS FULL
;
DEFINE FIELD text ON page TYPE string ASSERT string::len($value) > 0
	PERMISSIONS FULL
;
DEFINE FIELD embedding ON page TYPE array<float> ASSERT array::len($value) = 1536
	PERMISSIONS FULL
;
DEFINE FIELD createdAt ON page VALUE time::now()
	PERMISSIONS FULL
;
DEFINE FIELD updatedAt ON page VALUE time::now()
	PERMISSIONS FULL
;
DEFINE INDEX textIndex ON page FIELDS text SEARCH ANALYZER blank_snowball_ger BM25(1.2,0.75) DOC_IDS_ORDER 100 DOC_LENGTHS_ORDER 100 POSTINGS_ORDER 100 TERMS_ORDER 100 DOC_IDS_CACHE 100 DOC_LENGTHS_CACHE 100 POSTINGS_CACHE 100 TERMS_CACHE 100 HIGHLIGHTS;
DEFINE INDEX mTreeEmbeddingCosineIndex ON page FIELDS embedding MTREE DIMENSION 1536 DIST COSINE TYPE F64 CAPACITY 40 DOC_IDS_ORDER 100 DOC_IDS_CACHE 100 MTREE_CACHE 100;


--- DOCUMENT

DEFINE TABLE document TYPE ANY SCHEMAFULL
	PERMISSIONS NONE
;
DEFINE FIELD filePath ON document TYPE string ASSERT string::len($value) > 0
	PERMISSIONS FULL
;
DEFINE FIELD pages ON document VALUE <future> {
	RETURN (SELECT * FROM page:[
		$parent.id,
		NONE
	]..[
		$parent.id
	]);
}
	PERMISSIONS FULL
;
DEFINE FIELD createdAt ON document VALUE time::now()
	PERMISSIONS FULL
;
DEFINE FIELD updatedAt ON document VALUE time::now()
	PERMISSIONS FULL
;
DEFINE INDEX uniqueFilePathIndex ON document FIELDS filePath UNIQUE;