import style from "./style.module.scss";
import { type FC } from "react";

export const Layout: FC = () => {
  return (
    <>
      <header className={style.header}></header>
      <main className={style.main}></main>
      <footer className={style.footer}></footer>
    </>
  );
};
