import React,{useState} from 'react';

import styles from '_styles/Practice.module.scss';

export const Practice = () => {
    const [darkmode, setdarkmode] = useState(false);
    const topclass= darkmode? `${styles.container} ${styles.dark}` : styles.container;

    return (
        <div className={topclass}>
            <div className={styles.card}>
                <p>Testing Darkmode</p>
                <button 
                    onClick={()=>setdarkmode(x=>!x)}
                >
                    Toggle Darkmode
                </button>
            </div>
        </div>
    );
};

export default Practice;