/**
 * Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import ExternalLayout from 'layouts/External';
import { Button } from 'components';
import { useTranslation } from 'react-i18next';
import useAuth from 'helpers/hooks/useAuth';
import Styled from './styled';
import { keycloakInstance } from 'config/keycloak';
import { useHistory } from 'react-router-dom';
import accountService from 'services/auth';
import { clearTokens, getAccessToken } from 'helpers/localStorage/tokens';
import { setCurrentUser } from 'helpers/localStorage/currentUser';
import useFlashMessage from 'helpers/hooks/useFlashMessage';

function KeycloakAuth() {
  const { t } = useTranslation();
  const { login } = useAuth();
  const history = useHistory();
  const { showErrorFlash } = useFlashMessage();

  keycloakInstance.onAuthSuccess = () => {
    const accessToken = getAccessToken() || keycloakInstance.token;

    accountService
      .createAccountFromKeycloak(accessToken)
      .then((result) => {
        const userData = result?.data?.content;

        setCurrentUser(userData);

        history.replace('/home');
      })
      .catch(() => {
        showErrorFlash(t('API_ERRORS.KEYCLOAK_LOGIN'));

        setTimeout(() => {
          clearTokens();
          keycloakInstance.logout();
        }, 3200);
      });
  };

  return (
    <ExternalLayout>
      <Styled.Wrapper>
        <Button
          width={200}
          tabIndex={0}
          rounded
          text={t('LOGIN_SCREEN.KEYCLOAK')}
          onClick={() => login()}
        />
      </Styled.Wrapper>
    </ExternalLayout>
  );
}

export default KeycloakAuth;
