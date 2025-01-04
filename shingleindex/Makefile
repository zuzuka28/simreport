startup_env:
	python3.11 -m venv ./.venv
	( \
       source ./.venv/bin/activate; \
       pip install -r requirements.txt; \
    )

shutdown_env:
	rm -rf ./.venv

run_api:
	( \
       source ./.venv/bin/activate; \
       python3.11 -m src.handler; \
    )

