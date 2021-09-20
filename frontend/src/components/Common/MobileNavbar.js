import React    from 'react';

import { Link } from 'react-router-dom';

import {  
    Search as SearchIcon, 
    ArrowDropDown as ArrowDropDownIcon 
} from '@material-ui/icons';

import styles from '_styles/NavbarMobile.module.css';

import profilePicPlaceholder from 'assets/avatar-placeholder.svg';

const MobileNavbar = ({open}) => {

    

    return (
        <div className={ open ? `${styles.root} ${styles.open}` : styles.root}>
            <section className={styles.section}>
                <div className={styles.pfpFlex}>
                    <Link
                      className={styles.pfpContainer} 
                      to={routes.getProfile( authenticationService.currentUserValue.id ) }
                    >
                      <img 
                        className={styles.pfp} 
                        src={profilePicPath || profilePicPlaceholder}
                      />
                    </Link>

                    <div 
                      className={styles.dropDownIconContainer}
                      onClick={handleClick}
                    >
                      <ArrowDropDownIcon 
                        className={styles.dropDownIcon}
                      />
                    </div>
                </div>
            </section>
            <section className={styles.section}>
                
            </section>
            <section className={styles.section}> 
            </section>

        </div>
    )
}

export default MobileNavbar
