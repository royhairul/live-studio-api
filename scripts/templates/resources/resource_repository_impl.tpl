import (
	"gorm.io/gorm"

	"{{Module}}/internal/domains/{{feature}}/entity"
)

type {{Feature}}RepositoryImpl struct {
	DB *gorm.DB
}

func New{{Feature}}Repository(db *gorm.DB) {{Feature}}Repository {
	return &{{Feature}}RepositoryImpl{DB: db}
}

// Create implements {{Feature}}Repository.
func (r *{{Feature}}RepositoryImpl) Create(data *entity.{{Feature}}) (*entity.{{Feature}}, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements {{Feature}}Repository.
func (r *{{Feature}}RepositoryImpl) FindAll() ([]*entity.{{Feature}}, error) {
	var items []*entity.{{Feature}}
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements {{Feature}}Repository.
func (r *{{Feature}}RepositoryImpl) FindByID(id string) (*entity.{{Feature}}, error) {
	var item entity.{{Feature}}
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements {{Feature}}Repository.
func (r *{{Feature}}RepositoryImpl) Update(data *entity.{{Feature}}) (*entity.{{Feature}}, error) {
	if err := r.DB.Model(&entity.{{Feature}}{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements {{Feature}}Repository.
func (r *{{Feature}}RepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.{{Feature}}{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
