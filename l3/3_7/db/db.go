package db

import (
	"fmt"
	"os"
	"warehouseControl/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB_URL")

	// Connect with logging enabled to see what's happening
	newLogger := logger.Default
	if os.Getenv("GIN_MODE") == "debug" {
		newLogger = newLogger.LogMode(logger.Info)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Printf("Failed to connect to DB: %v\n", err)
		panic("Failed to connect to DB")
	}

	// Migrate the schema - this will create tables but not the trigger
	DB.AutoMigrate(&models.Inventory{}, &models.HistoryRecord{})

	// Create the trigger for history tracking
	err = createHistoryTrigger()
	if err != nil {
		fmt.Printf("Failed to create history trigger: %v\n", err)
	}
}

func createHistoryTrigger() error {
	// Create the function for the trigger - using GORM table names (pluralized by GORM)
	functionSQL := `
	CREATE OR REPLACE FUNCTION inventory_history_trigger_func()
	RETURNS TRIGGER AS $$
	BEGIN
		IF TG_OP = 'INSERT' THEN
			INSERT INTO history_records (item_id, action, new_values, changed_by, changed_at)
			VALUES (NEW.id, 'INSERT', row_to_json(NEW), COALESCE(current_user, 'system'), CURRENT_TIMESTAMP);
			RETURN NEW;
		ELSIF TG_OP = 'UPDATE' THEN
			INSERT INTO history_records (item_id, action, old_values, new_values, changed_by, changed_at)
			VALUES (OLD.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), COALESCE(current_user, 'system'), CURRENT_TIMESTAMP);
			RETURN NEW;
		ELSIF TG_OP = 'DELETE' THEN
			INSERT INTO history_records (item_id, action, old_values, changed_by, changed_at)
			VALUES (OLD.id, 'DELETE', row_to_json(OLD), COALESCE(current_user, 'system'), CURRENT_TIMESTAMP);
			RETURN OLD;
		END IF;
		RETURN NULL; -- Result is ignored since this is an AFTER trigger
	END;
	$$ LANGUAGE plpgsql;
	`

	// Create the trigger - GORM uses plural table names
	triggerSQL := `
	DROP TRIGGER IF EXISTS inventory_history_trigger ON inventories;
	CREATE TRIGGER inventory_history_trigger
		AFTER INSERT OR UPDATE OR DELETE ON inventories
		FOR EACH ROW EXECUTE FUNCTION inventory_history_trigger_func();
	`

	// Execute the SQL statements
	result := DB.Exec(functionSQL)
	if result.Error != nil {
		return result.Error
	}

	result = DB.Exec(triggerSQL)
	return result.Error
}
