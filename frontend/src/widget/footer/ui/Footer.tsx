import styles from "./Footer.module.scss";
import logo from "assets/icons/pepicons-print_hamburger-circle-filled_monochrome.svg";
import { navbarLinks, mediaLinks } from "shared/config";
import { IconLink } from "shared/ui";

export const Footer = () => {
  return (
    <div className={styles.wrapper}>
      <div className={styles.content}>
        <div className={styles.logoContainer}>
          <img className={styles.logo} src={logo} alt={"logo"} />
          <div className={styles.titleContainer}>
            <p className={styles.logoTitle}>food</p>
            <p className={styles.logoTitle}>adviser</p>
          </div>
        </div>
        <div className={styles.navbarContainer}>
          <NavBar items={navbarLinks} theme={"dark"} />
        </div>
        <ul className={styles.mediaContainer}>
          {mediaLinks.map((media) => (
            <li key={media.link} className={styles.mediaItem}>
              <IconLink
                icon={media.icon}
                link={media.link}
                className={styles.mediaIcon}
                target={"_blank"}
              />
            </li>
          ))}
        </ul>
      </div>
      <div className={styles.copyrightContainer}>
        <p className={styles.copyright}>
          Copyright: &copy;{new Date().getFullYear()} Myself.
        </p>
      </div>
    </div>
  );
};
