package analyzer

import (
	"database/sql"
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"time"
	"unsafe"

	_ "github.com/mattn/go-sqlite3"
)

// Path to the database file
const FXParamDBFile = "reaper_fx_params.db"

// WriteFXParamsToDB analyzes FX parameters and writes them directly to a SQLite database
func WriteFXParamsToDB() error {
	// Lock the current goroutine to the OS thread for UI operations
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("FX Parameter DB Writer started")

	// Get selected track
	track, err := reaper.GetSelectedTrack()
	if err != nil {
		return fmt.Errorf("please select a track with FX to analyze: %v", err)
	}

	// Get track info
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		return fmt.Errorf("error getting track info: %v", err)
	}

	// Check if track has FX
	if trackInfo.NumFX == 0 {
		return fmt.Errorf("selected track has no FX, please add FX to analyze")
	}

	// Show confirmation
	confirm, err := reaper.YesNoBox(
		fmt.Sprintf("This will analyze all parameters on %d FX on track '%s' and store them in a database.\n\nThis may take some time for complex plugins.\n\nProceed?",
			trackInfo.NumFX, trackInfo.Name),
		"FX Parameter DB Writer")

	if err != nil || !confirm {
		logger.Info("User cancelled parameter database writing")
		return nil
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Start timing
	startTime := time.Now()

	// For each FX on the track
	for fxIndex := 0; fxIndex < trackInfo.NumFX; fxIndex++ {
		// Get FX name
		fxName, err := reaper.GetTrackFXName(track, fxIndex)
		if err != nil {
			logger.Error("Failed to get FX name for index %d: %v", fxIndex, err)
			continue
		}

		logger.Info("Processing FX #%d: %s", fxIndex+1, fxName)

		// Insert FX and get ID
		fxID, err := getOrCreateFX(db, fxName)
		if err != nil {
			logger.Error("Failed to save FX to database: %v", err)
			continue
		}

		// Get parameter count
		paramCount, err := reaper.GetTrackFXParamCount(track, fxIndex)
		if err != nil {
			logger.Error("Failed to get parameter count for FX #%d: %v", fxIndex+1, err)
			continue
		}

		// For each parameter
		for paramIndex := 0; paramIndex < paramCount; paramIndex++ {
			// Process parameter
			err := processParameter(db, track, fxIndex, paramIndex, fxID)
			if err != nil {
				logger.Error("Error processing parameter %d for FX %s: %v", paramIndex, fxName, err)
				continue
			}
		}
	}

	// Calculate duration
	duration := time.Since(startTime)
	logger.Info("Database writing complete in %v", duration.Round(time.Millisecond))

	// Show completion message
	reaper.MessageBox(
		fmt.Sprintf("Analysis complete! All FX parameters stored in database in %v.",
			duration.Round(time.Millisecond)),
		"FX Parameter DB Writer")

	return nil
}

// Initialize the SQLite database
func initDatabase() (*sql.DB, error) {
	dbPath := "/Users/conor/Dev/external/reaper-go-extension/fx-dump.db"
	logger.Info("Opening database at: %s", dbPath)

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create tables if they don't exist
	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

// Create database tables
func createTables(db *sql.DB) error {
	// Use transaction for better performance and safety
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create FX table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS fx (
		fx_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		return err
	}

	// Create parameter table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS parameter (
		param_id INTEGER PRIMARY KEY AUTOINCREMENT,
		fx_id INTEGER NOT NULL,
		param_index INTEGER NOT NULL,
		name TEXT NOT NULL,
		is_toggle BOOLEAN,
		normal_step REAL,
		small_step REAL,
		large_step REAL,
		min_formatted TEXT,
		max_formatted TEXT,
		FOREIGN KEY (fx_id) REFERENCES fx(fx_id) ON DELETE CASCADE,
		UNIQUE(fx_id, param_index)
	)`)
	if err != nil {
		return err
	}

	// Create parameter sample table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS parameter_sample (
		sample_id INTEGER PRIMARY KEY AUTOINCREMENT,
		param_id INTEGER NOT NULL,
		normalized_value REAL NOT NULL,
		formatted_value TEXT NOT NULL,
		FOREIGN KEY (param_id) REFERENCES parameter(param_id) ON DELETE CASCADE,
		UNIQUE(param_id, normalized_value)
	)`)
	if err != nil {
		return err
	}

	// Create indexes
	_, err = tx.Exec(`
	CREATE INDEX IF NOT EXISTS idx_fx_name ON fx(name);
	CREATE INDEX IF NOT EXISTS idx_param_fx_id ON parameter(fx_id);
	CREATE INDEX IF NOT EXISTS idx_sample_param_id ON parameter_sample(param_id);
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get or create FX record and return its ID
func getOrCreateFX(db *sql.DB, fxName string) (int64, error) {
	// Check if FX already exists
	var fxID int64
	err := db.QueryRow("SELECT fx_id FROM fx WHERE name = ?", fxName).Scan(&fxID)
	if err == nil {
		// FX exists, return ID
		return fxID, nil
	}

	// Insert new FX
	result, err := db.Exec("INSERT INTO fx (name) VALUES (?)", fxName)
	if err != nil {
		return 0, err
	}

	// Get the inserted ID
	return result.LastInsertId()
}

// Process a single parameter and store it in the database
func processParameter(db *sql.DB, track unsafe.Pointer, fxIndex, paramIndex int, fxID int64) error {
	// Get parameter name
	paramName, err := reaper.GetTrackFXParamName(track, fxIndex, paramIndex)
	if err != nil {
		return fmt.Errorf("failed to get parameter name: %v", err)
	}

	// Get parameter range
	_, min, max, err := reaper.GetTrackFXParamValueWithRange(track, fxIndex, paramIndex)
	if err != nil {
		return fmt.Errorf("failed to get parameter range: %v", err)
	}

	// Get min/max formatted values
	minFormatted, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, min)
	if err != nil {
		minFormatted = ""
	}

	maxFormatted, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, max)
	if err != nil {
		maxFormatted = ""
	}

	// Get step sizes and toggle status
	var step, smallStep, largeStep float64
	var isToggle bool
	success := reaper.TrackFX_GetParameterStepSizes(track, fxIndex, paramIndex, &step, &smallStep, &largeStep, &isToggle)
	if !success {
		// Default values if API call fails
		step, smallStep, largeStep = 0, 0, 0
		isToggle = false
	}

	// Start transaction for better performance
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update parameter
	var paramID int64
	err = tx.QueryRow(`
		SELECT param_id FROM parameter WHERE fx_id = ? AND param_index = ?
	`, fxID, paramIndex).Scan(&paramID)

	if err == nil {
		// Parameter exists, update it
		_, err = tx.Exec(`
			UPDATE parameter SET 
				name = ?,
				is_toggle = ?,
				normal_step = ?,
				small_step = ?,
				large_step = ?,
				min_formatted = ?,
				max_formatted = ?
			WHERE param_id = ?
		`, paramName, isToggle, step, smallStep, largeStep, minFormatted, maxFormatted, paramID)
		if err != nil {
			return err
		}
	} else {
		// Insert new parameter
		result, err := tx.Exec(`
			INSERT INTO parameter (
				fx_id, param_index, name, is_toggle, normal_step, small_step, large_step,
				min_formatted, max_formatted
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, fxID, paramIndex, paramName, isToggle, step, smallStep, largeStep, minFormatted, maxFormatted)
		if err != nil {
			return err
		}

		paramID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	// Generate sample points based on the parameter's characteristics
	var samplePoints []float64

	// Always start with 0.0
	samplePoints = append(samplePoints, 0.0)

	// Use parameter's exact smallStep value when available
	if isToggle {
		// For toggle parameters, we only need 0.0 and 1.0
		// Already added 0.0, just need to add 1.0 below
	} else if smallStep > 0.0 {
		// Parameter has defined steps - use the exact smallStep without any limits
		// This may generate very large numbers of samples for parameters with tiny step sizes
		for point := smallStep; point < 1.0; point += smallStep {
			samplePoints = append(samplePoints, point)
		}

		logger.Info("    Sampling parameter at exact smallStep=%.6f (%d points)",
			smallStep, len(samplePoints)+1) // +1 for the final 1.0 we'll add
	} else {
		// Parameter has undefined or zero step size
		// Use a reasonable distribution with more points in important ranges
		samplePoints = append(samplePoints,
			0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09,
			0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45,
			0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 0.85, 0.9, 0.95,
			0.96, 0.97, 0.98, 0.99)

		logger.Info("    Sampling parameter with default distribution (smallStep=0)")
	}

	// Always end with exactly 1.0
	samplePoints = append(samplePoints, 1.0)

	// Process each sample point
	for _, point := range samplePoints {
		// Get formatted value at this normalized value
		formattedValue, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, point)
		if err != nil {
			logger.Warning("Failed to get formatted value for point %.2f: %v", point, err)
			continue
		}

		// Insert or update sample
		_, err = tx.Exec(`
			INSERT OR REPLACE INTO parameter_sample (
				param_id, normalized_value, formatted_value
			) VALUES (?, ?, ?)
		`, paramID, point, formattedValue)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Register the FX parameter DB writer action
func RegisterFXParamDBWriter() error {
	actionID, err := reaper.RegisterMainAction("GO_PARAM_DB_WRITER", "Go: Save FX Parameters to Database")
	if err != nil {
		return fmt.Errorf("failed to register FX parameter DB writer action: %v", err)
	}

	logger.Info("FX Parameter DB Writer registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_PARAM_DB_WRITER", handleFXParamDBWriter)
	return nil
}

// Handler for the FX parameter DB writer action
func handleFXParamDBWriter() {
	err := WriteFXParamsToDB()
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error: %v", err), "FX Parameter DB Writer")
	}
}
