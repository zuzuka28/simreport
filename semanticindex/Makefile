startup_env:
	python3.11 -m venv ./.venv
	( \
       source ./.venv/bin/activate; \
       pip install -r requirements.txt; \
       python -m spacy download ru_core_news_sm; \
    )

shutdown_env:
	rm -rf ./.venv

run_api:
	( \
       source ./.venv/bin/activate; \
       python3.11 -m src.service_handler; \
    )

run_indexer:
	( \
       source ./.venv/bin/activate; \
       python3.11 -m src.indexer_handler; \
    )
