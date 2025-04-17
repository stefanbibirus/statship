# Retro Synth Relationship

Aplicație web cu temă retro pentru vizualizarea apropierii și depărtării într-o relație dintre doi utilizatori, folosind o animație cu dublă elice (double helix).

## Caracteristici

- **Autentificare securizată** cu JWT și bcrypt
- **Sistem de invitație** cu coduri unice pentru formarea relațiilor
- **Animație double helix** care reflectă vizual apropierea și distanța dintre parteneri
- **Actualizări în timp real** folosind WebSockets
- **Design retro** cu efecte CRT, culori neon și fonturi pixelate

## Tehnologii

### Frontend
- React.js
- React Router pentru navigare
- HTML5 Canvas pentru animație
- Comunicare API cu Axios
- WebSockets pentru actualizări în timp real

### Backend
- Go (Golang) cu Fiber framework
- PostgreSQL pentru baza de date
- JWT pentru autentificare
- WebSockets pentru comunicare în timp real

## Structura proiectului

### Frontend (React)
```
relationship-helix-frontend/
├── public/
├── src/
│   ├── components/ (componente UI)
│   ├── context/ (gestionare stare globală)
│   ├── services/ (comunicare cu API)
│   └── styles/ (CSS)
```

### Backend (Go)
```
relationship-helix-backend/
├── cmd/
│   └── server/ (punctul de intrare)
├── internal/
│   ├── api/ (handlere, middleware și rute)
│   ├── config/ (configurație aplicație)
│   ├── db/ (acces bază de date și migrări)
│   ├── models/ (structuri de date)
│   └── utils/ (funcții utilitare)
```

## Instalare și rulare

### Cerințe
- Node.js (16+)
- Go (1.18+)
- PostgreSQL
- Git

### Pași pentru rulare locală

#### Backend
1. Clonează repository-ul
2. Navighează în directorul backend
   ```bash
   cd relationship-helix-backend
   ```
3. Instalează dependențele Go
   ```bash
   go mod download
   ```
4. Configurează baza de date PostgreSQL
   ```sql
   CREATE DATABASE relationship_helix;
   ```
5. Setează variabilele de mediu în fișierul `.env` (vezi `.env.example`)
6. Rulează aplicația
   ```bash
   go run cmd/server/main.go
   ```

#### Frontend
1. Navighează în directorul frontend
   ```bash
   cd relationship-helix-frontend
   ```
2. Instalează dependențele Node.js
   ```bash
   npm install
   ```
3. Setează variabilele de mediu în fișierul `.env` (vezi `.env.example`)
4. Rulează aplicația
   ```bash
   npm start
   ```
5. Accesează aplicația în browser la `http://localhost:3000`

## Deploy

### Backend (Railway sau Render)
1. Creează un cont pe Railway sau Render
2. Conectează repository-ul GitHub
3. Configurează variabilele de mediu (vezi `.env.example`)
4. Adaugă un serviciu PostgreSQL
5. Implementează aplicația

### Frontend (GitHub Pages)
1. Actualizează `package.json` pentru GitHub Pages
2. Configurează variabila `REACT_APP_API_URL` pentru URL-ul backend-ului
3. Rulează `npm run deploy`

## Utilizare

1. Creează un cont nou sau autentifică-te
2. Generează un cod de invitație și trimite-l partenerului tău
3. Partenerul tău introduce codul pentru a stabili relația
4. Vizualizați animația double helix care reprezintă relația voastră
5. Actualizați-vă poziția (apropiat/distant) și urmăriți în timp real schimbările

## Mockup-uri

- [Login Screen](/mockups/login-mockup.svg)
- [Dashboard](/mockups/dashboard-mockup.svg)
- [Invite Code](/mockups/invite-code-mockup.svg)
- [Settings](/mockups/settings-mockup.svg)

## Licență

Acest proiect este disponibil sub licența MIT. Consultă fișierul `LICENSE` pentru detalii.

## Autori

Dezvoltat ca parte a unui proiect demonstrativ pentru vizualizarea stării relațiilor interpersonale într-un mod creativ și estetic.