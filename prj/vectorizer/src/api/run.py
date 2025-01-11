import os
import click
from src.container import Container
from src.api.preprocess.routes import preprocess_router
import uvicorn
from fastapi import FastAPI
from fastapi.middleware import Middleware
from fastapi.middleware.cors import CORSMiddleware

from src.config import config


def init_routers(app_: FastAPI) -> None:
    app_.include_router(preprocess_router)


def make_middleware() -> list[Middleware]:
    middleware = [
        Middleware(
            CORSMiddleware,
            allow_origins=["*"],
            allow_credentials=True,
            allow_methods=["*"],
            allow_headers=["*"],
        ),
    ]
    return middleware


def create_app() -> FastAPI:
    _ = Container()

    app_ = FastAPI(
        title="SimRep Preprocessor",
        description="SimRep Preprocessor",
        version="1.0.0",
        middleware=make_middleware(),
    )
    init_routers(app_=app_)
    return app_


app = create_app()


@click.command()
@click.option(
    "--env",
    type=click.Choice(["local", "dev", "prod"], case_sensitive=False),
    default="local",
)
@click.option(
    "--debug",
    type=click.BOOL,
    is_flag=True,
    default=False,
)
def main(env: str, debug: bool):
    os.environ["ENV"] = env
    os.environ["DEBUG"] = str(debug)

    uvicorn.run(
        app="src.api.run:app",
        host=config.APP_HOST,
        port=config.APP_PORT,
        reload=True if config.ENV != "production" else False,
        workers=1,
    )


if __name__ == "__main__":
    main()
