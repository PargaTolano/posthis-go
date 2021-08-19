import React,{ Children, cloneElement } from 'react';

import styles from '_styles/DialogFollow.module.css';

export const DialogFollow = ({children, open, title, onClose}) => {

  const childrenWithProps = Children.map(children, (child, index) => {
    return cloneElement(child, {
      onClose,
    });
  });

  return (
    <div className={open? `${styles.container} ${styles.open}` : styles.container}>

      <div 
        className={styles.darkener}
        onClick={onClose}
      ></div>
      <div className={styles.dialog}>

        <div className={styles.dialogTitle}>
          <h3 className={styles.title}>
              {title}
              <span className={styles.close}></span>
          </h3>
        </div>
          {childrenWithProps}
      </div>
    </div>
  );
};

export default DialogFollow;